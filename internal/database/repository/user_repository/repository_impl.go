package user_repository

import (
	"errors"

	"github.com/vwency/microservices_golang/internal/database/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) RunMigrations() error {
	if err := r.db.AutoMigrate(&models.User{}); err != nil {
		return status.Errorf(codes.Internal, "failed to run migrations: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, status.Error(codes.NotFound, "user not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
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
		return nil, status.Error(codes.InvalidArgument, "username and email cannot both be empty")
	}

	err := query.First(&user).Error
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, status.Error(codes.NotFound, "user not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to retasdsadrieve user: %v", err)
	}
}

func (r *UserRepositoryImpl) AddUser(user *models.User) error {
	err := r.db.Create(user).Error
	switch {
	case err == nil:
		return nil
	case isDuplicateError(err):
		return status.Error(codes.AlreadyExists, "user already exists")
	default:
		return status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
}

func (r *UserRepositoryImpl) UpdateUserTokens(userID, hashedRefreshToken, hashedAccessToken string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"hashed_refresh_token": hashedRefreshToken,
			"hashed_access_token":  hashedAccessToken,
		})

	if result.Error != nil {
		return status.Errorf(codes.Internal, "failed to update tokens: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return status.Error(codes.NotFound, "user not found")
	}
	return nil
}

// Helper function to check for duplicate key errors
func isDuplicateError(err error) bool {
	switch err {
	case gorm.ErrDuplicatedKey:
		return true
	// Handle other database specific duplicate errors if needed
	default:
		return false
	}
}
