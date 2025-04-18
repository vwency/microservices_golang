package user_repository

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) RunMigrations() error {
	return RunUserMigrations(r.db)
}

func (r *UserRepositoryImpl) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	query := r.db.Model(&models.User{})

	switch {
	case username != "" && email != "":
		query = query.Where("username = ? OR email = ?", username, email)
	case username != "":
		query = query.Where("username = ?", username)
	case email != "":
		query = query.Where("email = ?", email)
	default:
		return nil, errors.New("username and email cannot both be empty")
	}

	result := query.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) AddUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepositoryImpl) UpdateUserTokens(userID, HashedRefreshToken, HashedAccessToken string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"hashed_refresh_token": HashedRefreshToken,
			"hashed_access_token":  HashedAccessToken,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update tokens: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
