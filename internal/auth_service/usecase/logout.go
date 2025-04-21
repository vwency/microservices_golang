package auth_service_usecase

import (
	"context"
	"fmt"
	"strings"

	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *AuthUsecase) Logout(ctx context.Context, username string, accessToken string) (bool, error) {
	uc.logger.Info("Attempting logout",
		zap.String("username", username),
		zap.String("ip", getIPFromContext(ctx)))

	if username == "" || accessToken == "" {
		uc.logger.Warn("Missing logout credentials",
			zap.String("username", username))
		return false, status.Error(codes.InvalidArgument, "username and access_token are required")
	}

	// Получить пользователя
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: username,
	})
	if err != nil {
		uc.logger.Error("Failed to get user for logout",
			zap.String("username", username),
			zap.Error(err))
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found during logout attempt",
			zap.String("username", username))
		return false, nil
	}

	// Проверка — пользователь уже вышел?
	if isTokenEmptyOrNone(getUserResp.HashedAccessToken) && isTokenEmptyOrNone(getUserResp.HashedRefreshToken) {
		uc.logger.Info("User already logged out",
			zap.String("user_id", getUserResp.UserId),
			zap.String("username", username))
		return true, nil
	}

	// Сравнение токена
	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, accessToken, getUserResp.HashedAccessToken)
	if err != nil {
		uc.logger.Error("Access token comparison failed",
			zap.String("username", username),
			zap.Error(err))
		return false, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
	}
	if !match {
		uc.logger.Warn("Invalid access token provided",
			zap.String("username", username))
		return false, status.Error(codes.Unauthenticated, "access token mismatch")
	}

	// Очистка токенов
	logoutRequest := &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: "none",
		HashedAccessToken:  "none",
	}

	_, err = uc.dbClient.UpdateUser(ctx, logoutRequest)
	if err != nil {
		uc.logger.Error("Failed to clear user tokens during logout",
			zap.String("user_id", getUserResp.UserId),
			zap.Error(err))

		if isTokenValidationError(err) {
			uc.logger.Warn("Retrying logout with placeholder tokens",
				zap.String("user_id", getUserResp.UserId))

			_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
				UserId:             getUserResp.UserId,
				HashedRefreshToken: "none",
				HashedAccessToken:  "none",
			})
			if err != nil {
				uc.logger.Error("Fallback logout failed",
					zap.String("user_id", getUserResp.UserId),
					zap.Error(err))
				return false, fmt.Errorf("fallback logout failed: %v", err)
			}
		} else {
			return false, fmt.Errorf("logout failed: %v", err)
		}
	}

	uc.logger.Info("User logged out successfully",
		zap.String("user_id", getUserResp.UserId),
		zap.String("username", username))

	return true, nil
}

func isTokenValidationError(err error) bool {
	if status, ok := status.FromError(err); ok {
		return status.Code() == codes.InvalidArgument && strings.Contains(status.Message(), "token fields")
	}
	return false
}

func isTokenEmptyOrNone(token string) bool {
	return token == "" || strings.ToLower(token) == "none"
}
