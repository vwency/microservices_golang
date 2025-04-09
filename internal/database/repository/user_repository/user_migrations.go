package user_repository

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

func RunUserMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
