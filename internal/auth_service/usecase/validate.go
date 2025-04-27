package auth_service_usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// ValidateAccessToken validates the access token.
func (uc *AuthUsecase) ValidateAccessToken(ctx context.Context, token string) (*ValidateResult, error) {
	// Token must be provided
	if token == "" {
		return nil, errors.New("access token is required")
	}

	// Validate token and extract claims
	claims, err := uc.jwtManager.ValidateToken(token)
	if err != nil {
		uc.logger.Error("Failed to validate access token",
			zap.Error(err),
			zap.String("token", token))
		return nil, fmt.Errorf("failed to validate access token: %w", err)
	}

	// Extract userID from claims
	userID, ok := claims["UserID"].(string)
	if !ok {
		uc.logger.Error("userID is missing or not a string",
			zap.Any("claims", claims))
		return nil, errors.New("invalid token: userID missing or not a string")
	}

	// Extract roles from claims and check their type
	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		uc.logger.Error("roles are missing or not an array",
			zap.Any("claims", claims))
		return nil, errors.New("invalid token: roles missing or not an array")
	}

	// Convert roles to []string
	var rolesStr []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			uc.logger.Error("invalid role type",
				zap.Any("role", role))
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		rolesStr = append(rolesStr, roleStr)
	}

	// Fetch user from the database using userID (changed from Username to UserId)
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		UserId: &userID, // Исправлено: используем UserId вместо Username
	})
	if err != nil {
		uc.logger.Error("Failed to fetch user during token validation",
			zap.Error(err),
			zap.String("userID", userID))
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Check if the user exists in the database
	if !getUserResp.Found {
		uc.logger.Error("User not found during token validation",
			zap.String("userID", userID))
		return nil, errors.New("user not found")
	}

	// Compare the token hash with the stored hashed access token in the database
	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		uc.logger.Error("Failed to compare token hashes",
			zap.Error(err),
			zap.String("userID", userID))
		return nil, fmt.Errorf("failed to compare token hashes: %w", err)
	}

	// Check if the token hash matches
	if !match {
		uc.logger.Warn("Access token hash mismatch",
			zap.String("userID", userID))
		return nil, errors.New("invalid access token")
	}

	// Verify the token's expiry time
	expiry, ok := claims["exp"].(float64)
	if !ok {
		uc.logger.Error("Token expiry missing or invalid",
			zap.Any("claims", claims))
		return nil, errors.New("invalid token: expiry missing or not a valid int64")
	}

	// Convert the expiry from float64 to int64
	expiryInt64 := int64(expiry)

	// Check if the token is expired
	if time.Now().Unix() > expiryInt64 {
		uc.logger.Warn("Access token is expired",
			zap.String("userID", userID),
			zap.Int64("expiry", expiryInt64))
		return nil, errors.New("access token has expired")
	}

	// Return the result of validation
	return &ValidateResult{
		Valid:     true,
		UserID:    userID,
		Roles:     rolesStr,
		ExpiresAt: expiryInt64,
	}, nil
}
