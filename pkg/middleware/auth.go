package middleware

import (
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
			log.WithField("message", "Missing \"Authorization\" header").Error("Missing \"Authorization\" header")
			h.Unauthorized(c, "Missing 'Authorization' header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ValidateToken(tokenString)

		if err != nil {
			log.WithField("error", err).Error("Failed validating token")
			h.Unauthorized(c, "Failed validating token")
			c.Abort()
			return
		}

		c.Set("UserID", claims.UserID.String())

		c.Next()
	}
}
