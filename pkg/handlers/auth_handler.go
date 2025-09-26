package handlers

import (
	"fmt"
	"share-docs/pkg/auth"
	"share-docs/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
	userService services.UserServiceInterface
}

type RegisterRequest struct {
	Email     string    `json:"email" binding:"required,email,max=255"`
	Password  string    `json:"password" binding:"required,min=8,max=128"`
	FirstName string    `json:"first_name" binding:"omitempty,max=100"`
	LastName  string    `json:"last_name" binding:"omitempty,max=100"`
	BirthDate time.Time `json:"birth_date" binding:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

func NewAuthHandler(userService services.UserServiceInterface, baseHandler BaseHandler) *AuthHandler {
	return &AuthHandler{
		BaseHandler: baseHandler,
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		h.BadRequest(c, fmt.Sprintf("Invalid request data: %v", err))
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, &req.BirthDate)

	if err != nil {
		fmt.Printf("Error: %v", err)
		switch err {
		case services.ErrEmailAlreadyExists:
			h.BadRequest(c, "Email already registered")
		case services.ErrInvalidEmail:
			h.BadRequest(c, "Invalid email format")
		case services.ErrWeakPassword:
			h.BadRequest(c, "Password does not meet security requirements")
		default:
			h.InternalError(c, "Failed to create account")
		}
		return
	}

	h.Created(c, user, "Account created successfully")
}

func (h *AuthHandler) Login(c *gin.Context) {
	log := h.GetLogger(c)

	var req LoginRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		h.BadRequest(c, fmt.Sprintf("Invalid request data: %v", err))
	}

	user, err := h.userService.GetUserByEmail(req.Email)

	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			h.BadRequest(c, "Invalid email or password")
		default:
			h.InternalError(c, "Failed to login user")
		}
		return
	}

	err = h.userService.ValidatePassword(user.Password, req.Password)
	if err != nil {
		h.InternalError(c, "Failed validating password")
		return
	}

	// Generate JWT token
	userID := user.ID
	userEmail := user.Email

	tokenPair, err := auth.GenerateTokenPair(userID, userEmail)

	if err != nil {
		log.WithFields(map[string]any{
			"user_id": userID,
			"email":   userEmail,
			"error":   err.Error(),
		}).Error("JWT failed signature!")
		h.InternalError(c, "Failed to generate JWT token")
		return
	}

	h.Success(c, tokenPair, "")
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	log := h.GetLogger(c)

	var req RefreshAccessTokenRequest

	if err := h.BindAndValidate(c, &req); err != nil {
		log.WithError(err).Error("Invalid request!")
		h.BadRequest(c, "Invalid request body! 'refresh_token' required")
		return
	}

	log.WithField("refresh_token", req.RefreshToken).Info("Token from request")

	rtc, err := auth.ValidateToken(req.RefreshToken, auth.RefreshToken)

	if err != nil {
		log.WithError(err).Error("Invalid refresh token!")
		h.Unauthorized(c, "Refresh token has expired! Login to generate a new one!")
		return
	}

	accessToken, err := auth.RefreshAccessToken(*rtc)

	if err != nil {
		log.WithError(err).Error("Failed to issue new access_token")
		h.InternalError(c, "Failed to issue new access_token")
		return
	}

	var res = &auth.TokenPair{
		AccessToken:  *accessToken,
		RefreshToken: req.RefreshToken,
	}

	h.Success(c, res, "Successfully refreshed token")
}
