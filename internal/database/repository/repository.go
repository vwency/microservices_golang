package repository

import (
	user_repository "github.com/vwency/microservices_golang/internal/database/repository/user_repository"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo user_repository.UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepo: user_repository.NewUserRepository(db),
	}
}
