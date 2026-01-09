package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jixlox0/studoto-backend/internal/errors"
	"github.com/jixlox0/studoto-backend/pkg/auth"
)

type AuthMiddleware struct {
	jwtAuth *auth.JWTAuth
}

func NewAuthMiddleware(jwtAuth *auth.JWTAuth) *AuthMiddleware {
	return &AuthMiddleware{jwtAuth: jwtAuth}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for authentication token in headers
		// Support both X-Auth-Key and x-auth-token headers
		authKey := c.GetHeader("X-Auth-Key")
		if authKey == "" {
			authKey = c.GetHeader("x-auth-token")
		}
		if authKey == "" {
			authKey = c.GetHeader("X-Auth-Token")
		}
		
		if authKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErrXAuthKeyRequired.Error(),
			})
			c.Abort()
			return
		}

		// Remove any whitespace
		token := strings.TrimSpace(authKey)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErrXAuthKeyEmpty.Error(),
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := m.jwtAuth.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErrInvalidToken.Error(),
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
