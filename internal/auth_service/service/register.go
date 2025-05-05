package service

import (
	"context"

	"github.com/go-kit/kit/log/level"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	ip := getIPFromContext(ctx)
	_ = level.Info(s.logger).Log(
		"msg", "Attempting registration",
		"username", req.Username,
		"ip", ip,
	)

	if req.Username == "" || req.Password == "" || req.Email == "" {
		_ = level.Warn(s.logger).Log(
			"msg", "Missing registration credentials",
			"username", req.Username,
		)
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}

	hashedPassword, err := authutils.GenHash(req.Username, req.Password, nil)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to hash password",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Generate tokens first
	roles := []interface{}{"user"}
	payload := map[string]interface{}{
		"Username": req.Username, // Use username as identifier
		"Roles":    roles,
	}

	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to generate access token",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, refreshExpiresAt, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to generate refresh token",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to hash access token",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to hash refresh token",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "failed to hash refresh token: %v", err)
	}

	// Create user with all fields including tokens
	addUserReq := &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	}

	addUserResp, err := s.dbClient.AddUser(ctx, addUserReq)
	if err != nil || !addUserResp.GetSuccess() {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to add user to user_database",
			"username", req.Username,
			"err", err,
			"db_message", addUserResp.GetMessage(),
		)
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", addUserResp.GetMessage())
	}

	// No need for separate update since we set tokens during creation
	_ = level.Info(s.logger).Log(
		"msg", "User registered successfully",
		"username", req.Username,
		"accessTokenExpiry", accessExpiresAt,
		"refreshTokenExpiry", refreshExpiresAt,
	)

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
