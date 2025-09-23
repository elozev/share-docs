package routes

import (
	"fmt"
	"net/http"
	"share-docs/pkg/db"
	"share-docs/pkg/handlers"
	"share-docs/pkg/logger"
	"share-docs/pkg/middleware"
	"share-docs/pkg/services"
	"share-docs/pkg/util"

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

	logConfig := logger.LogConfig{
		Level:       util.GetEnv("LOG_LEVEL", "info"),
		Environment: util.GetEnv("ENVIRONMENT", "development"),
		OutputPath:  util.GetEnv("LOG_OUTPUT", "stdout"),
		ServiceName: "share-docs",
		Version:     "1.0.0",
	}

	log, err := logger.NewLogger(logConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialise logger: %v", err))
	}

	r.Use(middleware.LoggingMiddleware(log))

	database := db.Connect()

	userService := services.NewUserService(database)
	baseHandler := handlers.NewBaseHandler(database, log)
	userHandler := handlers.NewUserHandler(userService, *baseHandler)

	api := r.Group("/api/v1")

	setupUserRoutes(api, userHandler)
	return r
}
