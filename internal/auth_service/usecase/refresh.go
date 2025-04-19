package auth_service_usecase

import (
	"context"
	"fmt"

	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"github.com/vwency/microservices_golang/utils/authutils"
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
	ip := getIPFromContext(ctx)
	uc.logger.Info("Attempting token refresh",
		zap.String("ip", ip))

	// Step 1: Validate the refresh token
	claims, err := uc.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token",
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Step 2: Get user from database
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

	// Step 3: Verify the refresh token matches the hashed version in DB
	match, err := authutils.ComparePasswordAndHash(tokenHashPepper, refreshToken, getUserResp.HashedRefreshToken)
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

	// Step 4: Generate new tokens
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(claims.UserID, claims.Roles)
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

	// Step 5: Hash the new refresh token before storing
	hashedRefreshToken, err := authutils.GenerateFromPassword(tokenHashPepper, newRefreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token",
			zap.String("user_id", claims.UserID),
			zap.Error(err),
			zap.String("ip", ip))
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	// Step 6: Update user in database with new tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             claims.UserID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  "", // Clear access token as it's short-lived
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
		zap.Time("access_token_expiry", expiresAt),
		zap.Time("refresh_token_expiry", refreshExpiresAt),
		zap.String("ip", ip))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
