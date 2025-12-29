// Package middleware
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/horaoen/go-backend-clean-architecture/internal/tokenutil"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "missing or invalid authorization header"})
			c.Abort()
			return
		}

		authToken := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("x-user-id", userID)
		c.Next()
	}
}
