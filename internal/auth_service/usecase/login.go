package auth_service_usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"github.com/vwency/microservices_golang/utils/authutils"
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

// контекстный ключ для IP
type contextKey string

const ipContextKey = contextKey("ip")

// getIPFromContext извлекает IP из контекста, добавленного middleware
func getIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipContextKey).(string); ok {
		return ip
	}
	return "unknown"
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	uc.logger.Info("Attempting login for user",
		zap.String("username", username),
		zap.String("ip", getIPFromContext(ctx)))

	// Validate input
	if username == "" || password == "" {
		uc.logger.Warn("Empty credentials provided",
			zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	// Get user from database
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: username,
	})
	if err != nil {
		uc.logger.Error("Database operation failed",
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found",
			zap.String("username", username))
		return nil, ErrUserNotFound
	}

	// Verify password with username-dependent hashing
	match, err := authutils.ComparePasswordAndHash(username, password, getUserResp.HashedPassword)
	if err != nil {
		uc.logger.Error("Password comparison failed",
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if !match {
		uc.logger.Warn("Invalid password provided",
			zap.String("username", username),
			zap.String("ip", getIPFromContext(ctx)))
		return nil, ErrInvalidCredentials
	}

	// Use default "user" role since roles aren't available in the response
	roles := []string{"user"}

	// Generate new tokens
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(getUserResp.Username, roles)
	if err != nil {
		uc.logger.Error("Access token generation failed",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	refreshToken, _, err := uc.jwtManager.GenerateRefreshToken(getUserResp.Username, roles)
	if err != nil {
		uc.logger.Error("Refresh token generation failed",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	// Hash tokens before storing
	hashedRefreshToken, err := authutils.GenerateFromPassword(tokenHashPepper, refreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	hashedAccessToken, err := authutils.GenerateFromPassword(tokenHashPepper, accessToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash access token",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	// Update user with new hashed tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	uc.logger.Info("Login successful",
		zap.String("user_id", getUserResp.UserId),
		zap.String("username", username),
		zap.Time("token_expiry", expiresAt))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
