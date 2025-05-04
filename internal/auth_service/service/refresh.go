package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log/level"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

// Service defines the interface for our auth service
type Service interface {
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
}

// TokenPair represents a pair of access and refresh tokens

func (s *service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Get the IP address from context
	ip := getIPFromContext(ctx)
	_ = level.Info(s.logger).Log("msg", "Attempting token refresh", "ip", ip)

	// Validate the refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Invalid refresh token", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract UserID and Roles from the claims map
	userID, ok := claims["UserID"].(string)
	if !ok {
		_ = level.Error(s.logger).Log("msg", "Invalid claim: UserID not found or not a string", "ip", ip)
		return nil, fmt.Errorf("invalid claim: UserID not found or not a string")
	}

	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		_ = level.Error(s.logger).Log("msg", "Invalid claim: Roles not found or not an array", "ip", ip)
		return nil, fmt.Errorf("invalid claim: Roles not found or not an array")
	}

	// Convert roles to []string
	var roles []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			_ = level.Error(s.logger).Log("msg", "Invalid role type", "role", role, "ip", ip)
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		roles = append(roles, roleStr)
	}

	// Fetch user details based on the user ID from the claims
	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: &userID})
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "UserDatabase operation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		_ = level.Warn(s.logger).Log("msg", "User not found", "user_id", userID, "ip", ip)
		return nil, ErrUserNotFound
	}

	// Compare the refresh token with the stored hashed refresh token
	match, err := authutils.ComparePasswordAndHash(s.tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Token comparison failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if !match {
		_ = level.Warn(s.logger).Log("msg", "Refresh token mismatch", "user_id", userID, "ip", ip)
		return nil, ErrInvalidToken
	}

	// Create payload map for token generation
	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	// Generate new access and refresh tokens
	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Access token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	newRefreshToken, refreshExpiresAt, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Refresh token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	// Hash the generated access and refresh tokens
	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash access token", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, newRefreshToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash refresh token", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	// Update the user's stored tokens in the database
	_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to update user tokens", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	_ = level.Info(s.logger).Log(
		"msg", "Tokens refreshed successfully",
		"user_id", userID,
		"access_token_expiry", accessExpiresAt,
		"refresh_token_expiry", refreshExpiresAt,
		"ip", ip,
	)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
)
