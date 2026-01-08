package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jixlox0/studoto-backend/pkg/auth"
	"github.com/jixlox0/studoto-backend/pkg/i18n"
)

type AuthMiddleware struct {
	jwtAuth *auth.JWTAuth
}

func NewAuthMiddleware(jwtAuth *auth.JWTAuth) *AuthMiddleware {
	return &AuthMiddleware{jwtAuth: jwtAuth}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.GetLanguageFromRequest(c)

		// Check for custom X-Auth-Key header
		authKey := c.GetHeader("X-Auth-Key")
		if authKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.T(lang, "x_auth_key_required"),
			})
			c.Abort()
			return
		}

		// Remove any whitespace
		token := strings.TrimSpace(authKey)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.T(lang, "x_auth_key_empty"),
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := m.jwtAuth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.T(lang, "invalid_token"),
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
