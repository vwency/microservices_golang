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
	// Step 1: Hash the password
	salt := []byte(username)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hashedPassword) == 0 {
		uc.logger.Error("Failed to hash password with argon2")
		return nil, fmt.Errorf("failed to hash password: %v", "hashing error")
	}

	encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword)

	// Step 2: Add user to the database (no user_id yet)
	addUserReq := &databasev1.AddUserRequest{
		Username:       username,
		HashedPassword: encodedPassword,
		Email:          email,
	}

	addUserResp, err := uc.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		uc.logger.Error("Failed to add user to database", zap.Error(err))
		return nil, err
	}

	if !addUserResp.Success {
		uc.logger.Error("Failed to add user to database", zap.String("message", addUserResp.Message))
		return nil, fmt.Errorf("failed to add user: %v", addUserResp.Message)
	}

	// Step 3: Retrieve user to get user_id
	getUserReq := &databasev1.GetUserRequest{
		Username: username,
	}
	getUserResp, err := uc.dbClient.GetUser(ctx, getUserReq)
	if err != nil {
		uc.logger.Error("Failed to get user from database", zap.Error(err))
		return nil, err
	}

	if !getUserResp.Found {
		uc.logger.Error("User not found after creation", zap.String("username", username))
		return nil, fmt.Errorf("user not found after creation")
	}

	// Step 4: Use user_id from the response to generate JWT tokens
	userID := getUserResp.Username // Assuming 'user_id' is returned in the 'username' field (modify if needed)

	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(userID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, _, err := uc.jwtManager.GenerateRefreshToken(userID, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	// Step 5: Return the response
	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
