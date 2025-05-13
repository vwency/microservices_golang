package service

import (
	"context"
	"fmt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

func (s *service) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	s.logger.Log("message", "Attempting login for user",
		"username", username,
		"ip", getIPFromContext(ctx))

	if username == "" || password == "" {
		s.logger.Log("message", "Empty credentials provided", "username", username)
		return nil, ErrInvalidCredentials
	}

	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &username,
	})
	if err != nil {
		s.logger.Log("message", "UserDatabase operation failed",
			"username", username,
			"error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		s.logger.Log("message", "User not found", "username", username)
		return nil, ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(username, password, getUserResp.HashedPassword)
	if err != nil {
		s.logger.Log("message", "Password comparison failed",
			"username", username,
			"error", err)
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		s.logger.Log("message", "Invalid password",
			"username", username,
			"ip", getIPFromContext(ctx))
		return nil, ErrInvalidCredentials
	}

	userID := getUserResp.UserId
	roles := []interface{}{"user"}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		s.logger.Log("message", "Failed to generate access token",
			"error", err,
			"user_id", userID)
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		s.logger.Log("message", "Failed to generate refresh token",
			"error", err,
			"user_id", userID)
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		s.logger.Log("message", "Failed to hash access token",
			"error", err,
			"user_id", userID)
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		s.logger.Log("message", "Failed to hash refresh token",
			"error", err,
			"user_id", userID)
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		s.logger.Log("message", "Failed to update user tokens",
			"user_id", userID,
			"error", err)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	s.logger.Log("message", "Login successful",
		"user_id", userID,
		"username", username,
		"token_expiry", accessExpiresAt)

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
