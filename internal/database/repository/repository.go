package repository

import (
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository" // Импортируем правильно
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo user_repository.UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	userRepo := user_repository.NewUserRepository(db)

	return &Repository{
		UserRepo: userRepo,
	}
}
