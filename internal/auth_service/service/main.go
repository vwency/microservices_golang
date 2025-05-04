package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
)

type service struct {
	dbClient    databasev1.DatabaseInitServiceClient
	jwtManager  *jwt.JWTManager
	logger      log.Logger
	tokenPepper string
}

func NewService(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger log.Logger,
	tokenPepper string,
) Service {
	return &service{
		dbClient:    dbClient,
		jwtManager:  jwtManager,
		logger:      logger,
		tokenPepper: tokenPepper,
	}
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error)
}
