package handlers

import (
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

func (h *UserHandler) GetUser(c *gin.Context) {
	log := h.GetLogger(c)
	userID := c.GetString("UserID")

	log.WithField("UserID", userID).Info("GetUser:userID")

	if userID == "" {
		h.InternalError(c, "UserID not found in context")
		return
	}

	user, err := h.userService.GetUserByID(userID)

	if err != nil {
		switch err {
		case services.ErrInvalidEmail:
			h.BadRequest(c, "Invalid email format")
		default:
			h.InternalError(c, "Failed to retrieve user by ID")
		}
		return
	}

	h.Success(c, user, "")
}
