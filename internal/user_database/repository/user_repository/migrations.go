package user_repository

import (
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"gorm.io/gorm"
)

func RunUserMigrations(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
