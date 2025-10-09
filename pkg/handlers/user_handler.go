package handlers

import (
	"share-docs/pkg/app/domain/userapp"
	"share-docs/pkg/services"

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

	appUser := userapp.ToAppUser(*user)
	h.Success(c, appUser, "")
}
