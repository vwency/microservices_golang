package auth_service_usecase

import (
	"context"

	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
)

type RefreshUsecase struct {
	dbClient   databasev1.DatabaseInitServiceClient
	jwtManager *jwt.JWTManager
	logger     *zap.Logger
}

func NewRefreshUsecase(dbClient databasev1.DatabaseInitServiceClient, jwtManager *jwt.JWTManager, logger *zap.Logger) *RefreshUsecase {
	return &RefreshUsecase{
		dbClient:   dbClient,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (uc *RefreshUsecase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := uc.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token", zap.Error(err))
		return nil, ErrInvalidToken
	}

	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: claims.UserID})
	if err != nil {
		uc.logger.Error("Failed to get user", zap.String("username", claims.UserID), zap.Error(err))
		return nil, err
	}

	if getUserResp.HashedRt != refreshToken {
		uc.logger.Warn("Refresh token mismatch", zap.String("username", claims.UserID))
		return nil, ErrInvalidToken
	}

	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.String("username", claims.UserID), zap.Error(err))
		return nil, err
	}

	newRefreshToken, _, err := uc.jwtManager.GenerateRefreshToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.String("username", claims.UserID), zap.Error(err))
		return nil, err
	}

	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		Username: claims.UserID,
		HashedRt: newRefreshToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user with refresh token", zap.String("username", claims.UserID), zap.Error(err))
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
