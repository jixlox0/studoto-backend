package main

import (
	"log"
	"os"

	"github.com/jixlox0/studoto-backend/internal/api"
	"github.com/jixlox0/studoto-backend/internal/config"
	"github.com/jixlox0/studoto-backend/internal/database"
	"github.com/jixlox0/studoto-backend/internal/middleware"
	"github.com/jixlox0/studoto-backend/internal/repository"
	"github.com/jixlox0/studoto-backend/internal/service"
	"github.com/jixlox0/studoto-backend/pkg/auth"
	"github.com/jixlox0/studoto-backend/pkg/i18n"
	"github.com/jixlox0/studoto-backend/pkg/oauth"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize i18n
	if err := i18n.Init(cfg.Server.DefaultLang); err != nil {
		log.Printf("Warning: Failed to initialize i18n: %v", err)
	}

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying sql.DB for proper connection closing
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	jwtAuth := auth.NewJWTAuth(cfg.JWT.SecretKey, cfg.JWT.ExpirationHours)
	oauthService := oauth.NewOAuthService(cfg.OAuth)
	userService := service.NewUserService(userRepo, jwtAuth)
	authService := service.NewAuthService(userRepo, jwtAuth, oauthService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtAuth)

	// Initialize handlers
	handlers := api.NewHandlers(userService, authService, authMiddleware)

	// Initialize router
	router := api.NewRouter(handlers, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
