package user_usecase

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"github.com/vwency/microservices_golang/internal/database/repository/user_repository"
	"go.uber.org/zap"
)

type UserUsecase struct {
	repo   user_repository.UserRepository
	logger *zap.Logger
}

func New(repo user_repository.UserRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *UserUsecase) GetUser(username, email string) (*models.User, error) {
	if username == "" && email == "" {
		return nil, errors.New("username and email cannot both be empty")
	}

	user, err := uc.repo.GetUserByUsernameOrEmail(username, email)
	if err != nil {
		uc.logger.Error("failed to get user",
			zap.String("username", username),
			zap.String("email", email),
			zap.Error(err))
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	return user, nil
}

func (uc *UserUsecase) CreateUser(params CreateUserParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("username", params.Username),
			zap.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	existingUser, err := uc.repo.GetUserByUsernameOrEmail(params.Username, params.Email)
	if err != nil {
		uc.logger.Error("failed to check user existence",
			zap.String("username", params.Username),
			zap.Error(err))
		return fmt.Errorf("check user existence failed: %w", err)
	}
	if existingUser != nil {
		uc.logger.Warn("user already exists",
			zap.String("username", params.Username))
		return ErrUserAlreadyExists
	}

	user := &models.User{
		Username: params.Username,
		Password: params.Password,
		HashedRt: params.HashedRt,
		HashedAt: params.HashedAt,
		Email:    &params.Email,
	}

	if err := uc.repo.AddUser(user); err != nil {
		uc.logger.Error("failed to create user",
			zap.String("username", params.Username),
			zap.Error(err))
		return fmt.Errorf("create user failed: %w", err)
	}

	uc.logger.Info("user created successfully",
		zap.String("username", params.Username))
	return nil
}

func (uc *UserUsecase) UpdateTokens(username, hashedRt, hashedAt string) error {
	if username == "" {
		uc.logger.Warn("empty username provided")
		return errors.New("username cannot be empty")
	}

	user, err := uc.GetUser(username, "")
	if err != nil {
		uc.logger.Error("failed to get user for token update",
			zap.String("username", username),
			zap.Error(err))
		return fmt.Errorf("get user failed: %w", err)
	}
	if user == nil {
		uc.logger.Warn("user not found for token update",
			zap.String("username", username))
		return ErrUserNotFound
	}

	if err := uc.repo.UpdateUserTokens(username, hashedRt, hashedAt); err != nil {
		uc.logger.Error("failed to update tokens",
			zap.String("username", username),
			zap.Error(err))
		return fmt.Errorf("update tokens failed: %w", err)
	}

	uc.logger.Info("tokens updated successfully",
		zap.String("username", username))
	return nil
}
