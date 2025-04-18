package auth_service_usecase

import (
	"context"
	"encoding/base64"

	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
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
	// Step 1: Validate the refresh token
	claims, err := uc.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token", zap.Error(err))
		return nil, ErrInvalidToken
	}

	// Step 2: Get user from database
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: claims.UserID})
	if err != nil {
		uc.logger.Error("Failed to get user",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	// Step 3: Verify the refresh token matches the hashed version in DB
	salt := []byte(claims.UserID)
	hashedIncomingToken := argon2.IDKey([]byte(refreshToken), salt, 1, 64*1024, 4, 32)
	encodedIncomingToken := base64.StdEncoding.EncodeToString(hashedIncomingToken)

	if getUserResp.HashedRefreshToken != encodedIncomingToken {
		uc.logger.Warn("Refresh token mismatch",
			zap.String("user_id", claims.UserID))
		return nil, ErrInvalidToken
	}

	// Step 4: Generate new tokens
	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	newRefreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(claims.UserID, claims.Roles)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	// Step 5: Hash the new refresh token before storing
	hashedNewRefreshToken := argon2.IDKey([]byte(newRefreshToken), salt, 1, 64*1024, 4, 32)
	encodedNewRefreshToken := base64.StdEncoding.EncodeToString(hashedNewRefreshToken)

	// Step 6: Update user in database with new tokens
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             claims.UserID,
		HashedRefreshToken: encodedNewRefreshToken,
		HashedAccessToken:  "", // Clear access token as it's short-lived
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Tokens refreshed successfully",
		zap.String("user_id", claims.UserID),
		zap.Time("access_token_expiry", expiresAt),
		zap.Time("refresh_token_expiry", refreshExpiresAt))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
