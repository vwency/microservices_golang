package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

// TokenPair represents the access and refresh token pair


// contextKey for extracting IP address
type contextKey string

const ipContextKey = contextKey("ip")



// LoginService defines the dependencies for login usecase
type LoginService struct {
	dbClient    databasev1.DatabaseInitServiceClient
	jwtManager  *jwt.JWTManager
	logger      *zap.Logger
	tokenPepper string
}

func NewLoginService(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
	tokenPepper string,
) *LoginService {
	return &LoginService{
		dbClient:    dbClient,
		jwtManager:  jwtManager,
		logger:      logger,
		tokenPepper: tokenPepper,
	}
}

// Login authenticates user and returns tokens
func (s *LoginService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	s.logger.Info("Attempting login for user", zap.String("username", username), zap.String("ip", getIPFromContext(ctx)))

	if username == "" || password == "" {
		s.logger.Warn("Empty credentials provided", zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &username,
	})
	if err != nil {
		s.logger.Error("UserDatabase operation failed", zap.String("username", username), zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		s.logger.Warn("User not found", zap.String("username", username))
		return nil, ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(username, password, getUserResp.HashedPassword)
	if err != nil {
		s.logger.Error("Password comparison failed", zap.String("username", username), zap.Error(err))
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		s.logger.Warn("Invalid password", zap.String("username", username), zap.String("ip", getIPFromContext(ctx)))
		return nil, ErrInvalidCredentials
	}

	userID := getUserResp.UserId
	roles := []interface{}{"user"}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	accessToken, accessExpiresAt, err := s.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		s.logger.Error("Failed to hash access token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}
	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		s.logger.Error("Failed to hash refresh token", zap.Error(err), zap.String("user_id", userID))
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		s.logger.Error("Failed to update user tokens", zap.String("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	s.logger.Info("Login successful", zap.String("user_id", userID), zap.String("username", username), zap.Time("token_expiry", accessExpiresAt))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}
