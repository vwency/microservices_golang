package auth_service_usecase

import (
	"context"
	"errors"

	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

type ValidateResult struct {
	Valid     bool
	UserID    string
	Roles     []string
	ExpiresAt int64
}

func (uc *AuthUsecase) ValidateAccessToken(ctx context.Context, token string) (*ValidateResult, error) {
	if token == "" {
		return nil, errors.New("access token is required")
	}

	claims, err := uc.jwtManager.ValidateToken(token)
	if err != nil {
		uc.logger.Error("Failed to validate access token", zap.Error(err))
		return nil, err
	}
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: claims.UserID})
	if err != nil {
		uc.logger.Error("Failed to fetch user during token validation", zap.Error(err), zap.String("userID", claims.UserID))
		return nil, err
	}
	if !getUserResp.Found {
		uc.logger.Error("User not found during token validation", zap.String("userID", claims.UserID))
		return nil, errors.New("user not found")
	}

	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		uc.logger.Error("Failed to compare token hashes", zap.Error(err), zap.String("userID", claims.UserID))
		return nil, err
	}
	if !match {
		uc.logger.Warn("Access token hash mismatch", zap.String("userID", claims.UserID))
		return nil, errors.New("invalid access token")
	}

	return &ValidateResult{
		Valid:     true,
		UserID:    claims.UserID,
		Roles:     claims.Roles,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}
