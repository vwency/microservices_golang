package service

import (
	"context"
	"errors"
	"fmt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func (s *service) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	refreshToken := req.GetRefreshToken()
	ip := getIPFromContext(ctx)
	s.logger.Log("msg", "Attempting token refresh", "ip", ip)

	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Log("msg", "Invalid refresh token", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	userID, ok := claims["UserID"].(string)
	if !ok {
		s.logger.Log("msg", "Invalid claim: UserID not found or not a string", "ip", ip)
		return nil, errors.New("invalid claim: UserID not found or not a string")
	}

	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		s.logger.Log("msg", "Invalid claim: Roles not found or not an array", "ip", ip)
		return nil, errors.New("invalid claim: Roles not found or not an array")
	}

	var roles []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			s.logger.Log("msg", "Invalid role type", "role", role, "ip", ip)
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		roles = append(roles, roleStr)
	}

	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		s.logger.Log("msg", "UserDatabase operation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		s.logger.Log("msg", "User not found", "user_id", userID, "ip", ip)
		return nil, ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(s.tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		s.logger.Log("msg", "Token comparison failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if !match {
		s.logger.Log("msg", "Refresh token mismatch", "user_id", userID, "ip", ip)
		return nil, ErrInvalidToken
	}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		s.logger.Log("msg", "Access token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	newRefreshToken, _, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		s.logger.Log("msg", "Refresh token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		s.logger.Log("msg", "Failed to hash access token", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, newRefreshToken, nil)
	if err != nil {
		s.logger.Log("msg", "Failed to hash refresh token", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		s.logger.Log("msg", "Failed to update user tokens", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	s.logger.Log(
		"msg", "Tokens refreshed successfully",
		"user_id", userID,
		"access_token_expiry", accessExpiresAt,
		"ip", ip,
	)

	return &authv1.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
