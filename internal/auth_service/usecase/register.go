package auth_service_usecase

import (
	"context"
	"fmt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) (*authv1.RegisterResponse, error) {
	if username == "" || password == "" || email == "" {
		return nil, fmt.Errorf("username, password and email are required")
	}

	hashedPassword, err := authutils.GenHash(username, password, nil)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	tempUserID := username

	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(tempUserID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err), zap.String("tempUserID", tempUserID))
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(tempUserID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.Error(err), zap.String("tempUserID", tempUserID))
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	hashedAccessToken, err := authutils.GenHash(uc.tokenPepper, accessToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash access token", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(uc.tokenPepper, refreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	addUserReq := &databasev1.AddUserRequest{
		Username:           username,
		HashedPassword:     hashedPassword,
		Email:              email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	}

	addUserResp, err := uc.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		uc.logger.Error("Failed to add user to database", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to add user: %v", err)
	}

	if !addUserResp.Success {
		uc.logger.Error("Database operation failed", zap.String("message", addUserResp.Message), zap.String("username", username))
		return nil, fmt.Errorf("database operation failed: %v", addUserResp.Message)
	}

	getUserReq := &databasev1.GetUserRequest{Username: username}
	getUserResp, err := uc.dbClient.GetUser(ctx, getUserReq)
	if err != nil {
		uc.logger.Error("Failed to retrieve user after creation", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to retrasddsasdieve user: %v", err)
	}

	if !getUserResp.Found {
		uc.logger.Error("User not found after creation", zap.String("username", username))
		return nil, fmt.Errorf("user not found after creation")
	}

	uc.logger.Info("User registered successfully",
		zap.String("userID", getUserResp.UserId),
		zap.String("username", username),
		zap.Time("accessTokenExpiry", expiresAt),
		zap.Time("refreshTokenExpiry", refreshExpiresAt))

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
