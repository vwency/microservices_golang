package auth_service_usecase

import (
	"context"
	"fmt"

	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

func (uc *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	ip := getIPFromContext(ctx)
	uc.logger.Info("Attempting token refresh",
		zap.String("ip", ip))

	claims, err := uc.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token",
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: claims.UserID})
	if err != nil {
		uc.logger.Error("Database operation failed",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found",
			zap.String("user_id", claims.UserID),
			zap.String("ip", ip))
		return nil, ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		uc.logger.Error("Token comparison failed",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if !match {
		uc.logger.Warn("Refresh token mismatch",
			zap.String("user_id", claims.UserID),
			zap.String("ip", ip))
		return nil, ErrInvalidToken
	}

	// Generate new tokens
	accessToken, accessExpiresAt, err := uc.jwtManager.GenerateAccessToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Access token generation failed",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	newRefreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Refresh token generation failed",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	hashedAccessToken, err := authutils.GenHash(uc.tokenPepper, accessToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash access token",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(uc.tokenPepper, newRefreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	uc.logger.Info("Tokens refreshed successfully",
		zap.String("user_id", claims.UserID),
		zap.Time("access_token_expiry", accessExpiresAt),
		zap.Time("refresh_token_expiry", refreshExpiresAt),
		zap.String("ip", ip))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}
