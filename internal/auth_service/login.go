package auth_service

import (
	"bytes"
	"context"
	"encoding/base64"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	getUserReq := &databasev1.GetUserRequest{
		Username: req.Username,
	}

	getUserResp, err := s.dbClient.GetUser(ctx, getUserReq)
	if err != nil {
		s.logger.Error("failed to get user from database", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if !getUserResp.Found {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	storedPassword, err := base64.StdEncoding.DecodeString(getUserResp.HashedPassword)
	if err != nil {
		s.logger.Error("failed to decode stored password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to decode password: %v", err)
	}

	salt := []byte(req.Username)
	hashedPassword := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)

	if !bytes.Equal(storedPassword, hashedPassword) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	roles := []string{"user"}
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(req.Username, roles)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(req.Username, roles)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	updateReq := &databasev1.UpdateUserRequest{
		Username: req.Username,
		HashedRt: refreshToken,
		AccessRt: "access-token",
	}

	_, err = s.dbClient.UpdateUser(ctx, updateReq)
	if err != nil {
		s.logger.Error("failed to update user refresh token", zap.Error(err))
	}

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
