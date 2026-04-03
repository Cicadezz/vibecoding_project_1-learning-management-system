package middleware

import (
	"net/http"
	"strings"

	"learning-growth-platform/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	svc *auth.Service
}

func NewAuthMiddleware(svc *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{svc: svc}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if strings.TrimSpace(header) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		userID, err := m.svc.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
