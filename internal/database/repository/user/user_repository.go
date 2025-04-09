package user

import (
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"github.com/vwency/microservices_golang/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewUserRepository(cfg config.ServiceConfig) (*userRepository, error) {
	dsn := cfg.Database.URL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &userRepository{db: db}, nil
}

func (r *userRepository) RunMigrations() error {
	if err := r.db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}
