package user_usecase

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/database/models"
	"go.uber.org/zap"
)

func (uc *UserUsecase) GetUser(params UserParams) (*models.User, error) {
	if params.Username == "" && params.Email == "" {
		return nil, errors.New("username and email cannot both be empty")
	}

	user, err := uc.repo.GetUserByUsernameOrEmail(params.Username, params.Email)
	if err != nil {
		uc.logger.Error("failed to get user",
			zap.String("username", params.Username),
			zap.String("email", params.Email),
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
		Username:           params.Username,
		HashedPassword:     params.HashedPassword,
		HashedRefreshToken: params.HashedRt,
		HashedAccessToken:  params.HashedAt,
		Email:              &params.Email,
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

func (uc *UserUsecase) UpdateTokens(params UpdateTokensParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("username", params.UserID),
			zap.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	user, err := uc.GetUser(UserParams{Username: params.UserID})
	if err != nil {
		uc.logger.Error("failed to get user for token update",
			zap.String("username", params.UserID),
			zap.Error(err))
		return fmt.Errorf("get user failed: %w", err)
	}
	if user == nil {
		uc.logger.Warn("user not found for token update",
			zap.String("username", params.UserID))
		return ErrUserNotFound
	}

	if err := uc.repo.UpdateUserTokens(params.UserID, params.HashedRefreshToken, params.HashedAccessToken); err != nil {
		uc.logger.Error("failed to update tokens",
			zap.String("username", params.UserID),
			zap.Error(err))
		return fmt.Errorf("update tokens failed: %w", err)
	}

	uc.logger.Info("tokens updated successfully",
		zap.String("username", params.UserID))
	return nil
}
