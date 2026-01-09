//go:build wireinject
// +build wireinject

// Package main provides dependency injection using Google Wire.
// This file is only compiled when the wireinject build tag is present.
// Run `make wire` or `go generate ./cmd/server` to generate wire_gen.go

package main

import (
	"github.com/google/wire"
	"github.com/jixlox0/studoto-backend/internal/api"
	"github.com/jixlox0/studoto-backend/internal/config"
	"github.com/jixlox0/studoto-backend/internal/database"
	"github.com/jixlox0/studoto-backend/internal/middleware"
	"github.com/jixlox0/studoto-backend/internal/repository"
	"github.com/jixlox0/studoto-backend/internal/service"
	"github.com/jixlox0/studoto-backend/pkg/auth"
	"github.com/jixlox0/studoto-backend/pkg/cache"
	"github.com/jixlox0/studoto-backend/pkg/oauth"
)

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

// InitializeApp initializes the application with all dependencies using Wire.
// This function is used by Wire to generate the dependency injection code.
// The generated code will be in wire_gen.go
func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// Config providers - extract fields from config struct
		provideDatabaseConfig,
		provideRedisConfig,
		provideJWTSecretKey,
		provideJWTExpirationHours,
		provideOAuthConfig,

		// Cache layer
		cache.NewRedisClient,
		cache.NewRedisCache,

		// Database layer
		database.NewConnection,

		// Repository layer
		repository.NewUserRepository,

		// Authentication & Authorization
		auth.NewJWTAuth,
		oauth.NewOAuthService,

		// Service layer
		service.NewUserService,
		service.NewAuthService,

		// Middleware
		middleware.NewAuthMiddleware,

		// API handlers
		api.NewHandlers,

		// Router
		api.NewRouter,

		// Application
		NewApp,
	)

	return nil, nil
}

// Config providers
// These functions extract specific fields from the config struct
// to provide them as individual dependencies to Wire.

// provideDatabaseConfig extracts the database configuration from the main config.
func provideDatabaseConfig(cfg *config.Config) config.DatabaseConfig {
	return cfg.Database
}

// provideJWTSecretKey extracts the JWT secret key from the config.
func provideJWTSecretKey(cfg *config.Config) string {
	return cfg.JWT.SecretKey
}

// provideJWTExpirationHours extracts the JWT expiration hours from the config.
func provideJWTExpirationHours(cfg *config.Config) int {
	return cfg.JWT.ExpirationHours
}

// provideOAuthConfig extracts the OAuth configuration from the main config.
func provideOAuthConfig(cfg *config.Config) config.OAuthConfig {
	return cfg.OAuth
}

// provideRedisConfig extracts the Redis configuration from the main config.
func provideRedisConfig(cfg *config.Config) cache.RedisConfig {
	return cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
}
