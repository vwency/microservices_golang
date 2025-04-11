package auth_service

import (
	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
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
		dbClient:   dbClient,
	}
}
