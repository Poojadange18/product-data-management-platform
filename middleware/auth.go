package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"product-data-management-platform/service"
	"strings"
)

func RequireAuth(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
			return
		}
		token, err := auth.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Next()
	}
}
