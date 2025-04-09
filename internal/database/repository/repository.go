package repository

import (
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"github.com/vwency/microservices_golang/internal/database/repository/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	RunMigrations() error
	GetUserByUsernameOrEmail(username, email string) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUserTokens(username, hashedRt, hashedAt string) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	repo := user.NewUserRepository(db)

	if err := repo.RunMigrations(); err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}

	return repo
}
