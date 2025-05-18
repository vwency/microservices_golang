package gokit

import (
	"github.com/go-kit/log"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
)

func NewAuthService(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger log.Logger,
	cfg config.ServiceConfig,
) service.AuthService {
	return service.NewService(dbClient, jwtManager, logger, cfg.Jwt.HashPepper)
}
