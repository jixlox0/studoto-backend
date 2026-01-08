package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jixlox0/studoto-backend/internal/config"
)

func NewRouter(handlers *Handlers, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.CORS.AllowedOrigins,
		AllowMethods:     cfg.Server.CORS.AllowedMethods,
		AllowHeaders:     cfg.Server.CORS.AllowedHeaders,
		ExposeHeaders:    cfg.Server.CORS.ExposedHeaders,
		AllowCredentials: cfg.Server.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.Server.CORS.MaxAge) * time.Second,
	}

	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/oauth/:provider", handlers.GetOAuthURL)
		auth.GET("/callback/:provider", handlers.OAuthCallback)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(handlers.authMiddleware.RequireAuth())
	{
		protected.GET("/profile", handlers.GetProfile)
	}

	return router
}

