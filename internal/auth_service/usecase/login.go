package auth_service_usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"time"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
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
		uc.logger.Error("Failed to get user",
			zap.String("username", username),
			zap.Error(err))
		return nil, err
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found", zap.String("username", username))
		return nil, ErrUserNotFound
	}

	// Verify password
	storedPassword, err := base64.StdEncoding.DecodeString(getUserResp.HashedPassword)
	if err != nil {
		uc.logger.Error("Failed to decode stored password",
			zap.String("username", username),
			zap.Error(err))
		return nil, err
	}

	salt := []byte(username)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if !bytes.Equal(storedPassword, hashedPassword) {
		uc.logger.Warn("Invalid credentials", zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	roles := []string{"user"} // Default role, adjust as needed
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(username, roles)
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			zap.String("username", username),
			zap.Error(err))
		return nil, ErrTokenGeneration
	}

	refreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(username, roles)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			zap.String("username", username),
			zap.Error(err))
		return nil, ErrTokenGeneration
	}

	// Hash tokens before storage
	hashedAccessToken := argon2.IDKey([]byte(accessToken), salt, 1, 64*1024, 4, 32)
	encodedAccessToken := base64.StdEncoding.EncodeToString(hashedAccessToken)

	hashedRefreshToken := argon2.IDKey([]byte(refreshToken), salt, 1, 64*1024, 4, 32)
	encodedRefreshToken := base64.StdEncoding.EncodeToString(hashedRefreshToken)

	// Update user with new hashed tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId, // Using UserId instead of Username
		HashedRefreshToken: encodedRefreshToken,
		HashedAccessToken:  encodedAccessToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens",
			zap.String("user_id", getUserResp.UserId),
			zap.String("username", username),
			zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Login successful",
		zap.String("user_id", getUserResp.UserId),
		zap.String("username", username),
		zap.Time("access_token_expiry", expiresAt),
		zap.Time("refresh_token_expiry", refreshExpiresAt))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
