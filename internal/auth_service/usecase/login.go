package auth_service_usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	uc.logger.Info("Attempting login for user", zap.String("username", username))

	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: username})
	if err != nil {
		uc.logger.Error("Failed to get user", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found", zap.String("username", username))
		return nil, ErrUserNotFound
	}

	storedPassword, err := base64.StdEncoding.DecodeString(getUserResp.HashedPassword)
	if err != nil {
		uc.logger.Error("Failed to decode stored password", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	salt := []byte(username)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if !bytes.Equal(storedPassword, hashedPassword) {
		uc.logger.Warn("Invalid credentials", zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	roles := []string{"user"}
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(username, roles)
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	refreshToken, _, err := uc.jwtManager.GenerateRefreshToken(username, roles)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		Username: username,
		HashedRt: refreshToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user with refresh token", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
