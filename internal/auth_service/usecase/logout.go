package auth_service_usecase

import (
	"context"

	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
)

type LogoutUsecase struct {
	dbClient   databasev1.DatabaseInitServiceClient
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
}

func NewLogoutUsecase(dbClient databasev1.DatabaseInitServiceClient, jwtManager *jwt.JWTManager, logger *zap.Logger) *LogoutUsecase {
	return &LogoutUsecase{
		dbClient:   dbClient,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (uc *LogoutUsecase) Logout(ctx context.Context, username string) (bool, error) {
	// First get the user to obtain their user_id
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: username})
	if err != nil {
		uc.logger.Error("Failed to get user", zap.String("username", username), zap.Error(err))
		return false, err
	}

	if !getUserResp.Found {
		uc.logger.Error("User not found", zap.String("username", username))
		return false, nil
	}

	// Update the user with empty refresh and access tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: "",
		HashedAccessToken:  "",
	})
	if err != nil {
		uc.logger.Error("Failed to update user during logout",
			zap.String("username", username),
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))
		return false, err
	}

	return true, nil
}
