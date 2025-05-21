package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
)

func getIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipContextKey).(string); ok {
		return ip
	}
	return "unknown"
}

func GetIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipContextKey).(string); ok {
		return ip
	}
	return "unknown"
}

func NewService(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger log.Logger,
	tokenPepper string,
) AuthService {
	return &service{
		dbClient:    dbClient,
		jwtManager:  jwtManager,
		logger:      logger,
		tokenPepper: tokenPepper,
	}
}
