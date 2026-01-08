package database

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jixlox0/studoto-backend/internal/config"
	"github.com/jixlox0/studoto-backend/internal/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var db *gorm.DB
	var err error
	maxRetries := 5
	retryDelay := 2 * time.Second

	// Retry connection with exponential backoff
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			log.Printf("Database connection attempt %d/%d failed, retrying in %v...", i+1, maxRetries, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	// Convert our migrations to gormigrate format
	migrationList := migrations.GetMigrations()
	gormigrateMigrations := make([]*gormigrate.Migration, len(migrationList))

	for i, mig := range migrationList {
		gormigrateMigrations[i] = &gormigrate.Migration{
			ID:       mig.ID,
			Migrate:  mig.Migrate,
			Rollback: mig.Rollback,
		}
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormigrateMigrations)

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
