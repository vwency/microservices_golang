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
	// Validate input
	if username == "" || password == "" || email == "" {
		return nil, fmt.Errorf("username, password and email are required")
	}

	// Step 1: Hash password with username-dependent hashinsg
	hashedPassword, err := authutils.GenerateFromPassword(username, password, nil)
	if err != nil {
		uc.logger.Error("Failed to hash password",
			zap.Error(err),
			zap.String("username", username))
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Step 2: Generate tokens
	tempUserID := username // Temporary ID, will be replaced with actual ID from DB

	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(tempUserID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			zap.Error(err),
			zap.String("tempUserID", tempUserID))
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(tempUserID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			zap.Error(err),
			zap.String("tempUserID", tempUserID))
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}
	addUserReq := &databasev1.AddUserRequest{
		Username:           username,
		HashedPassword:     hashedPassword,
		Email:              email,
		HashedAccessToken:  accessToken,  // Note: Consider hashing these too
		HashedRefreshToken: refreshToken, // Note: Consider hashing these too
	}

	addUserResp, err := uc.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		uc.logger.Error("Failed to add user to database",
			zap.Error(err),
			zap.String("username", username))
		return nil, fmt.Errorf("failed to add user: %v", err)
	}

	if !addUserResp.Success {
		uc.logger.Error("Database operation failed",
			zap.String("message", addUserResp.Message),
			zap.String("username", username))
		return nil, fmt.Errorf("database operation failed: %v", addUserResp.Message)
	}

	// Step 4: Get user details with actual user_id
	getUserReq := &databasev1.GetUserRequest{
		Username: username,
	}
	getUserResp, err := uc.dbClient.GetUser(ctx, getUserReq)
	if err != nil {
		uc.logger.Error("Failed to retrieve user after creation",
			zap.Error(err),
			zap.String("username", username))
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	if !getUserResp.Found {
		uc.logger.Error("User not found after creation",
			zap.String("username", username))
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
