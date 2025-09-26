package middleware

import (
	"fmt"
	"share-docs/pkg/auth"
	"share-docs/pkg/handlers"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(h handlers.BaseHandlerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.GetLogger(c)

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.WithField("message", "Missing \"Authorization\" header").Error("missing \"Authorization\" header")
			h.Unauthorized(c, "missing 'Authorization' header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ValidateToken(tokenString)

		if err != nil {
			log.WithField("error", err).Error("failed validating token")
			h.Unauthorized(c, fmt.Sprintf("%s", err.Error()))
			c.Abort()
			return
		}

		c.Set("UserID", claims.UserID.String())

		c.Next()
	}
}
