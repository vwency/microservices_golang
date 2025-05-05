package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	ip := getIPFromContext(ctx)
	_ = level.Info(s.logger).Log("msg", "Attempting to register user", "ip", ip, "username", req.Username)

	hashedPassword, err := authutils.GenHash(s.tokenPepper, req.Password, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Password hashing failed", "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	addUserResp, err := s.dbClient.AddUser(ctx, &databasev1.AddUserRequest{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	})
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "User creation failed", "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to add user: %w", err)
	}

	payload := map[string]interface{}{
		"UserID": addUserResp.UserId,
		"Roles":  []string{"user"},
	}

	accessToken, _, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Access token generation failed", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Refresh token generation failed", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash access token", "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash refresh token", "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             addUserResp.UserId,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	})
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to update user tokens", "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to update user tokens: %w", err)
	}

	_ = level.Info(s.logger).Log(
		"msg", "User registered successfully",
		"user_id", addUserResp.UserId,
		"ip", ip,
	)

	return &authv1.RegisterResponse{
		UserId:       addUserResp.UserId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 1).Unix(),
	}, nil
}
