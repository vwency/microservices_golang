package auth_service_usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
)

var (
	ErrTokenGeneration = errors.New("failed to generate tokens")
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	uc.logger.Info("Attempting login for user", zap.String("username", username))

	// Validate input
	if username == "" || password == "" {
		uc.logger.Warn("Empty username or password provided")
		return nil, ErrInvalidCredentials
	}

	// Get user from database
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: username})
	if err != nil {
		uc.logger.Error("Failed to get user from database",
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found in database", zap.String("username", username))
		return nil, ErrUserNotFound
	}

	// Verify password (in real app use proper hashing)
	if getUserResp.HashedPassword != password {
		uc.logger.Warn("Invalid password provided", zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	roles := []string{"user"} // Default role
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(getUserResp.UserId, roles)
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			zap.String("user_id", getUserResp.UserId),
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	refreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(getUserResp.UserId, roles)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			zap.String("user_id", getUserResp.UserId),
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}
	fmt.Printf(refreshExpiresAt.String())

	// Update user with new tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: refreshToken,
		HashedAccessToken:  accessToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens in database",
			zap.String("user_id", getUserResp.UserId),
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	uc.logger.Info("Login successful",
		zap.String("user_id", getUserResp.UserId),
		zap.String("username", username))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
