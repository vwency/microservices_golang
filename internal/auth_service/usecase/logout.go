package auth_service_usecase

import (
	"context"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
)

func (uc *AuthUsecase) Logout(ctx context.Context, username string) (bool, error) {
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: username})
	if err != nil {
		uc.logger.Error("Failed to get user",
			zap.String("username", username),
			zap.Error(err))
		return false, err
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found",
			zap.String("username", username))
		return false, nil
	}

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

	uc.logger.Info("User logged out successfully",
		zap.String("user_id", getUserResp.UserId),
		zap.String("username", username))

	return true, nil
}
