package handlers

import (
	"fmt"
	"share-docs/pkg/auth"
	"share-docs/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BaseHandler
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface, baseHandler BaseHandler) *UserHandler {
	return &UserHandler{
		BaseHandler: baseHandler,
		userService: userService,
	}
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

func (h *UserHandler) Register(c *gin.Context) {
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

func (h *UserHandler) Login(c *gin.Context) {
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

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.userService.GetUserByEmail(email)

	if err != nil {
		switch err {
		case services.ErrInvalidEmail:
			h.BadRequest(c, "Invalid email format")
		default:
			h.InternalError(c, "Failed to retrieve user by email")
		}
		return
	}

	h.Success(c, user, "")
}
