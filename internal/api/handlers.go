package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jixlox0/studoto-backend/internal/errors"
	"github.com/jixlox0/studoto-backend/internal/middleware"
	"github.com/jixlox0/studoto-backend/internal/models"
	"github.com/jixlox0/studoto-backend/internal/service"
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
func (h *Handlers) Signup(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.authService.Signup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) Signin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.authService.Signin(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetOAuthURL(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidProvider.Error()})
		return
	}

	url, err := h.authService.GetOAuthURL(provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handlers) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrCodeRequired.Error()})
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrNotAuthenticated.Error(),
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
			"error": errors.ErrInvalidUserIDType.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByID(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.ErrUserNotFound.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(user))
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
