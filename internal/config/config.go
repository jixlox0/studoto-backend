package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

type OAuthConfig struct {
	Google      GoogleOAuthConfig
	GitHub      GitHubOAuthConfig
	RedirectURL string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port string
	CORS CORSConfig
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func Load() (*Config, error) {
	expirationHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "studoto"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt(getEnv("REDIS_DB", "0"), 0),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationHours: expirationHours,
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			},
			GitHub: GitHubOAuthConfig{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			},
			RedirectURL: getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			CORS: CORSConfig{
				AllowedOrigins:   parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "*")),
				AllowedMethods:   parseStringSlice(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")),
				AllowedHeaders:   parseStringSlice(getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Auth-Key,x-auth-token,X-Auth-Token")),
				ExposedHeaders:   parseStringSlice(getEnv("CORS_EXPOSED_HEADERS", "Content-Length")),
				AllowCredentials: getEnv("CORS_ALLOW_CREDENTIALS", "true") == "true",
				MaxAge:           parseInt(getEnv("CORS_MAX_AGE", "86400"), 86400),
			},
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}
