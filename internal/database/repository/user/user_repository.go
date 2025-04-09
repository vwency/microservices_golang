package user

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// RunMigrations выполняет миграции для пользователей
func (r *userRepository) RunMigrations() error {
	if err := r.db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}
