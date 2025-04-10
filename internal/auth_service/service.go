package auth_service

import (
	"context"
	"fmt"

	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	authv1.UnimplementedAuthServiceServer
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
	dbClient   databasev1.DatabaseInitServiceClient
}

func NewAuthService(jwtManager *jwt.JWTManager, logger *zap.Logger, dbClient databasev1.DatabaseInitServiceClient) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
		logger:     logger.With(zap.String("service", "auth_service")),
		dbClient:   dbClient, // Инициализация gRPC клиента
	}
}

func (s *AuthService) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username, password, and email are required")
	}

	salt := []byte(req.Username)
	hashedPassword := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)

	if len(hashedPassword) == 0 {
		s.logger.Error("failed to hash password with argon2")
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", fmt.Errorf("hashing error"))
	}

	addUserReq := &databasev1.AddUserRequest{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		HashedRt:       "qweqwe",
		AccessRt:       "access-token",
		Email:          req.Email,
	}

	addUserResp, err := s.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		s.logger.Error("failed to add user to database", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", err)
	}

	if !addUserResp.Success {
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", addUserResp.Message)
	}

	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(req.Username, []string{"user"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(req.Username, []string{"user"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
