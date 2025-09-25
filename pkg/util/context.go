package util

import (
	"fmt"
	"share-docs/pkg/logger"

	"github.com/gin-gonic/gin"
)

func GetLoggerFromContext(c *gin.Context) *logger.Logger {
	if loggerInterface, exists := c.Get("Logger"); exists {
		if logger, ok := loggerInterface.(*logger.Logger); ok {
			return logger
		}
	}
	fmt.Println("Logger doesn't exists in context, creating one!")
	defaultLogger, _ := logger.NewLogger(logger.LogConfig{})
	return defaultLogger
}
