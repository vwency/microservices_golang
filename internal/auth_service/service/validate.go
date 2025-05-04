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

// ValidateResult contains the results of token validation.
type ValidateResult struct {
	Valid     bool
	UserID    string
	Roles     []string
	ExpiresAt int64
}

// ValidateService defines methods related to token validation.
type ValidateService interface {
	ValidateAccessToken(ctx context.Context, token string) (*ValidateResult, error)
}

type validateService struct {
	dbClient    databasev1.DatabaseInitServiceClient
	jwtManager  *jwt.JWTManager
	logger      *zap.Logger
	tokenPepper string
}

// NewValidateService creates a new ValidateService.
func NewValidateService(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
	tokenPepper string,
) ValidateService {
	return &validateService{
		dbClient:    dbClient,
		jwtManager:  jwtManager,
		logger:      logger,
		tokenPepper: tokenPepper,
	}
}

// ValidateAccessToken validates the access token.
func (s *validateService) ValidateAccessToken(ctx context.Context, token string) (*ValidateResult, error) {
	if token == "" {
		return nil, errors.New("access token is required")
	}

	// Validate JWT token
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		s.logger.Error("Failed to validate access token", zap.Error(err), zap.String("token", token))
		return nil, fmt.Errorf("failed to validate access token: %w", err)
	}

	// Extract userID and roles from the claims
	userID, ok := claims["UserID"].(string)
	if !ok {
		s.logger.Error("userID missing or invalid", zap.Any("claims", claims))
		return nil, errors.New("invalid token: userID missing or not a string")
	}

	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		s.logger.Error("roles missing or invalid", zap.Any("claims", claims))
		return nil, errors.New("invalid token: roles missing or not an array")
	}

	var rolesStr []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			s.logger.Error("invalid role type", zap.Any("role", role))
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		rolesStr = append(rolesStr, roleStr)
	}

	// Check if user exists
	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		UserId: &userID,
	})
	if err != nil {
		s.logger.Error("Failed to fetch user", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !getUserResp.Found {
		s.logger.Error("User not found", zap.String("userID", userID))
		return nil, errors.New("user not found")
	}

	// Compare token hashes for security
	match, err := authutils.ComparePasswordAndHash(s.tokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		s.logger.Error("Failed to compare token hashes", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to compare token hashes: %w", err)
	}

	if !match {
		s.logger.Warn("Access token hash mismatch", zap.String("userID", userID))
		return nil, errors.New("invalid access token")
	}

	// Validate expiration
	expiry, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Error("Token expiry missing or invalid", zap.Any("claims", claims))
		return nil, errors.New("invalid token: expiry missing or invalid")
	}

	expiryInt64 := int64(expiry)
	if time.Now().Unix() > expiryInt64 {
		s.logger.Warn("Access token expired", zap.String("userID", userID), zap.Int64("expiry", expiryInt64))
		return nil, errors.New("access token has expired")
	}

	return &ValidateResult{
		Valid:     true,
		UserID:    userID,
		Roles:     rolesStr,
		ExpiresAt: expiryInt64,
	}, nil
}
