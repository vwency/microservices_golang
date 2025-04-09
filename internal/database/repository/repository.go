package repository

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	RunMigrations() error
	GetUserByUsernameOrEmail(username, email string) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUserTokens(username, hashedRt, hashedAt string) error
}

// UserRepository struct implements the UserRepository interface
type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) RunMigrations() error {
	return r.db.AutoMigrate(&models.User{})
}

func (r *UserRepositoryImpl) GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImpl) AddUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepositoryImpl) UpdateUserTokens(username, hashedRt, hashedAt string) error {
	result := r.db.Model(&models.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"hashed_rt": hashedRt,
			"hashed_at": hashedAt,
		})
	return result.Error
}
