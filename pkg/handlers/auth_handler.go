package handlers

import (
	"fmt"
	"share-docs/pkg/auth"
	"share-docs/pkg/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
	userService services.UserServiceInterface
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
