package middleware

import (
	"share-docs/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		reqLogger := log.WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		c.Set("Logger", reqLogger)

		reqLogger.Info("Request started")

		c.Next()

		duration := time.Since(start)

		reqLogger.WithFields(map[string]interface{}{
			"status":      c.Writer.Status(),
			"size":        c.Writer.Size(),
			"duration":    duration.String(),
			"duration_ms": duration.Microseconds(),
		})

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				reqLogger.WithError(err.Err).Error("Request error occurred")
			}
		}
	}
}

func GetLoggerFromContext(c *gin.Context) *logger.Logger {
	if loggerInterface, exists := c.Get("logger"); exists {
		if logger, ok := loggerInterface.(*logger.Logger); ok {
			return logger
		}
	}
	defaultLogger, _ := logger.NewLogger(logger.LogConfig{})
	return defaultLogger
}
