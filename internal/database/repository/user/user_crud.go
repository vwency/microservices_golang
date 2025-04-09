package user

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AddUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserTokens обновляет токены пользователя (refresh_token и access_token)
func (r *userRepository) UpdateUserTokens(username, hashedRt, hashedAt string) error {
	if err := r.db.Model(&models.User{}).Where("username = ?", username).
		Update("refresh_token", hashedRt).
		Update("access_token", hashedAt).Error; err != nil {
		return err
	}
	return nil
}
