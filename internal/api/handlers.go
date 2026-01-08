package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jixlox0/studoto-backend/internal/middleware"
	"github.com/jixlox0/studoto-backend/internal/models"
	"github.com/jixlox0/studoto-backend/internal/service"
	"github.com/jixlox0/studoto-backend/pkg/i18n"
)

type Handlers struct {
	userService    service.UserService
	authService    service.AuthService
	authMiddleware *middleware.AuthMiddleware
}

func NewHandlers(userService service.UserService, authService service.AuthService, authMiddleware *middleware.AuthMiddleware) *Handlers {
	return &Handlers{
		userService:    userService,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

// Auth handlers
func (h *Handlers) Register(c *gin.Context) {
	lang := i18n.GetLanguageFromRequest(c)

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.TWithDefault(lang, "validation_error", err.Error()),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		errorMsg := i18n.TWithDefault(lang, "user_already_exists", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) Login(c *gin.Context) {
	lang := i18n.GetLanguageFromRequest(c)

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.TWithDefault(lang, "validation_error", err.Error()),
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		errorMsg := i18n.TWithDefault(lang, "invalid_credentials", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetOAuthURL(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	url, err := h.authService.GetOAuthURL(provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handlers) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	// Validate state if needed
	_ = state

	response, err := h.authService.OAuthLogin(provider, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// User handlers
func (h *Handlers) GetProfile(c *gin.Context) {
	lang := i18n.GetLanguageFromRequest(c)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": i18n.T(lang, "user_not_authenticated"),
		})
		return
	}

	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int:
		userIDUint = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.T(lang, "invalid_user_id_type"),
		})
		return
	}

	user, err := h.userService.GetUserByID(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": i18n.T(lang, "user_not_found"),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
