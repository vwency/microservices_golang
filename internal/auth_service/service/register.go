package service

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
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
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}

	// Генерируем user_id один раз
	userID := uuid.New().String()

	hashedPassword, err := authutils.GenHash(req.Username, req.Password, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash password", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Генерируем токены с user_id
	payload := map[string]interface{}{
		"UserID": userID, // Используем сгенерированный ID
		"Roles":  []interface{}{"user"},
	}

	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to generate access token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to generate refresh token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	// Хешируем токены
	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash access token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		_ = level.Error(s.logger).Log("msg", "Failed to hash refresh token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to hash refresh token: %v", err)
	}

	// Создаем запрос с user_id
	addUserReq := &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
		UserId:             &userID, // Используем тот же ID
	}

	addUserResp, err := s.dbClient.AddUser(ctx, addUserReq)
	if err != nil || !addUserResp.GetSuccess() {
		_ = level.Error(s.logger).Log("msg", "Failed to add user", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", err)
	}

	_ = level.Info(s.logger).Log(
		"msg", "User registered successfully",
		"user_id", userID,
		"username", req.Username,
	)

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
