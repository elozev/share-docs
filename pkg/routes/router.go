package routes

import (
	"net/http"
	"share-docs/pkg/db"
	"share-docs/pkg/handlers"
	"share-docs/pkg/services"

	"github.com/gin-gonic/gin"
)

func setupUserRoutes(r *gin.RouterGroup, userHandler *handlers.UserHandler) {

	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
	}

	user := r.Group("/user")
	{
		user.GET("/email/:email", userHandler.GetUserByEmail)
	}

	// TODO: add users group
}

// SetupRouter configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	database := db.Connect()

	userService := services.NewUserService(database)
	userHandler := handlers.NewUserHandler(userService)

	api := r.Group("/api/v1")

	setupUserRoutes(api, userHandler)
	return r
}
