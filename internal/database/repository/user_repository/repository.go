package user_repository

import "github.com/vwency/microservices_golang/internal/database/models"

type UserRepository interface {
	RunMigrations() error
	GetUserByID(id string) (*models.User, error)
	GetUserByUsernameOrEmail(username, email string) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUserTokens(userID, hashedRefreshToken, hashedAccessToken string) error
}
