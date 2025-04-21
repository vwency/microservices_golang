package auth_service_usecase

import (
	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"go.uber.org/zap"
)

type AuthUsecase struct {
	dbClient    databasev1.DatabaseInitServiceClient
	jwtManager  *jwt.JWTManager
	logger      *zap.Logger
	tokenPepper string
}

func NewAuthUsecase(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
	tokenPepper string,
) *AuthUsecase {
	return &AuthUsecase{
		dbClient:    dbClient,
		jwtManager:  jwtManager,
		logger:      logger,
		tokenPepper: tokenPepper,
	}
}
