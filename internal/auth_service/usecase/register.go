package auth_service_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) (*authv1.RegisterResponse, error) {
	if username == "" || password == "" || email == "" {
		return nil, fmt.Errorf("username, password and email are required")
	}

	// Generate a unique user_id
	userID := uuid.New().String()

	// Hash the password
	hashedPassword, err := authutils.GenHash(username, password, nil)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create the payload for token generation with the user_id
	roles := []interface{}{"user"} // Define roles for the user
	payload := map[string]interface{}{
		"UserID": userID, // Use user_id in the token payload
		"Roles":  roles,
	}

	// Generate access token and refresh token
	accessToken, accessExpiresAt, err := uc.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	// Hash the generated tokens
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

	// Create request to add user to the database with user_id
	addUserReq := &databasev1.AddUserRequest{
		Username:           username,
		HashedPassword:     hashedPassword,
		Email:              email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	}

	// Add the user to the database
	addUserResp, err := uc.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		uc.logger.Error("Failed to add user to user_database", zap.Error(err), zap.String("username", username))
		return nil, fmt.Errorf("failed to add user: %v", err)
	}

	if !addUserResp.Success {
		uc.logger.Error("UserDatabase operation failed", zap.String("message", addUserResp.Message), zap.String("username", username))
		return nil, fmt.Errorf("database_user operation failed: %v", addUserResp.Message)
	}

	// Update the user with the generated user_id
	updateUserReq := &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	}
	_, err = uc.dbClient.UpdateUser(ctx, updateUserReq)
	if err != nil {
		uc.logger.Error("Failed to update user with user_id", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to update user with user_id: %v", err)
	}

	// Log successful user registration
	uc.logger.Info("User registered successfully", zap.String("userID", userID), zap.String("username", username), zap.Time("accessTokenExpiry", accessExpiresAt), zap.Time("refreshTokenExpiry", refreshExpiresAt))

	// Return the response with the tokens
	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
