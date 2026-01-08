package migrations

import (
	"github.com/jixlox0/studoto-backend/internal/models"
	"github.com/jixlox0/studoto-backend/pkg/uuid"
	"gorm.io/gorm"
)

// GetMigrations returns all database migrations
func GetMigrations() []*Migration {
	return []*Migration{
		{
			ID: "20240101000001",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.User{})
			},
		},
		{
			ID: "20240101000002",
			Migrate: func(tx *gorm.DB) error {
				// Add UUID column if it doesn't exist
				if !tx.Migrator().HasColumn(&models.User{}, "uuid") {
					if err := tx.Migrator().AddColumn(&models.User{}, "uuid"); err != nil {
						return err
					}
					// Generate UUIDs for existing users that don't have one
					var users []models.User
					if err := tx.Where("uuid = '' OR uuid IS NULL").Find(&users).Error; err != nil {
						return err
					}
					for _, user := range users {
						userUUID := uuid.Generate(uuid.PrefixUser)
						if err := tx.Model(&user).Update("uuid", userUUID).Error; err != nil {
							return err
						}
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if tx.Migrator().HasColumn(&models.User{}, "uuid") {
					return tx.Migrator().DropColumn(&models.User{}, "uuid")
				}
				return nil
			},
		},
		// Add more migrations here as needed
	}
}

// Migration represents a database migration
type Migration struct {
	ID       string
	Migrate  func(*gorm.DB) error
	Rollback func(*gorm.DB) error
}
