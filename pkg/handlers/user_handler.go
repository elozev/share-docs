package handlers

import (
	"fmt"
	"share-docs/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BaseHandler
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
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
