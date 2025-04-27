package auth_service_usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	// Assuming the JWTManager is part of this package
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

var (
	ErrTokenGeneration = errors.New("failed to generate tokens")
)

// TokenPair represents the access and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// contextKey for extracting IP address
type contextKey string

const ipContextKey = contextKey("ip")

// getIPFromContext extracts IP address from the context (e.g., added by middleware)
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

	// Fetch user data from the database
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &username,
	})
	if err != nil {
		uc.logger.Error("UserDatabase operation failed",
			zap.String("username", username),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found",
			zap.String("username", username))
		return nil, ErrUserNotFound
	}

	// Verify password with the hashed password from the database
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

	// Extract user_id and roles from the database response (default "user" role)
	userID := getUserResp.UserId   // Use user_id from the database
	roles := []interface{}{"user"} // Modify if roles are fetched from the database

	// Create the payload for token generation, with user_id instead of username
	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	// Generate access token and refresh token using JWTManager
	accessToken, accessExpiresAt, err := uc.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, _, err := uc.jwtManager.GenerateRefreshToken(payload) // We don't need refreshExpiresAt
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	// Hash the generated tokens
	hashedAccessToken, err := authutils.GenHash(uc.tokenPepper, accessToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash access token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(uc.tokenPepper, refreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	// Update user record in the database with hashed tokens
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
		zap.Time("token_expiry", accessExpiresAt))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}
