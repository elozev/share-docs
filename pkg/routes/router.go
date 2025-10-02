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

func setupAuthRoutes(r *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}
}

func setupUserRoutes(r *gin.RouterGroup, userHandler *handlers.UserHandler) {
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(userHandler))
	{
		user.GET("/", userHandler.GetUser)
	}
}

func setupDocumentRoutes(r *gin.RouterGroup, documentHandler *handlers.DocHandler) {
	docs := r.Group("/docs")
	docs.Use(middleware.AuthMiddleware(documentHandler))
	{
		docs.POST("/", documentHandler.CreateDocument)
	}
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
	docService := services.NewDocumentService(database)
	storageType := util.MustGetEnv("STORAGE_TYPE")
	storageService := services.NewStorageService(storageType, log)

	baseHandler := handlers.NewBaseHandler(database, log)
	userHandler := handlers.NewUserHandler(userService, *baseHandler)
	authHandler := handlers.NewAuthHandler(userService, *baseHandler)
	docHandler := handlers.NewDocHandler(*docService, *storageService, *baseHandler)

	api := r.Group("/api/v1")

	setupAuthRoutes(api, authHandler)
	setupUserRoutes(api, userHandler)
	setupDocumentRoutes(api, docHandler)
	return r
}
