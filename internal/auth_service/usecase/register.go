package auth_service_usecase

import (
	"context"
	"encoding/base64"
	"fmt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) (*authv1.RegisterResponse, error) {
	// Validate input
	if username == "" || password == "" || email == "" {
		return nil, fmt.Errorf("username, password and email are required")
	}

	// Step 1: Hash the password with Argon2
	salt := []byte(username) // Using username as salt (consider random salt in production)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword)

	// Step 2: Generate tokens first (we'll need them for storage)
	// Temporary userID for token generation (will be replaced with actual ID from DB)
	tempUserID := username

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

	// Hash the tokens before storage
	hashedAccessToken := argon2.IDKey([]byte(accessToken), salt, 1, 64*1024, 4, 32)
	encodedAccessToken := base64.StdEncoding.EncodeToString(hashedAccessToken)

	hashedRefreshToken := argon2.IDKey([]byte(refreshToken), salt, 1, 64*1024, 4, 32)
	encodedRefreshToken := base64.StdEncoding.EncodeToString(hashedRefreshToken)

	// Step 3: Add user to database with hashed tokens
	addUserReq := &databasev1.AddUserRequest{
		Username:           username,
		HashedPassword:     encodedPassword,
		Email:              email,
		HashedAccessToken:  encodedAccessToken,
		HashedRefreshToken: encodedRefreshToken,
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

	// If we need to update tokens with actual user_id, we would do it here
	// But typically the tokens are valid with username as identifier

	uc.logger.Info("User registered successfully",
		zap.String("userID", getUserResp.UserId),
		zap.String("username", username),
		zap.Time("accessTokenExpiry", expiresAt),
		zap.Time("refreshTokenExpiry", refreshExpiresAt))

	// Return the original (unhashed) tokens to the client
	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
