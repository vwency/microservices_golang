package auth_service_usecase

import (
	"context"
	"fmt"

	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

func (uc *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Get the IP address from context
	ip := getIPFromContext(ctx)
	uc.logger.Info("Attempting token refresh", zap.String("ip", ip))

	// Validate the refresh token
	claims, err := uc.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token", zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract UserID and Roles from the claims map
	userID, ok := claims["UserID"].(string)
	if !ok {
		uc.logger.Error("Invalid claim: UserID not found or not a string", zap.String("ip", ip))
		return nil, fmt.Errorf("invalid claim: UserID not found or not a string")
	}

	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		uc.logger.Error("Invalid claim: Roles not found or not an array", zap.String("ip", ip))
		return nil, fmt.Errorf("invalid claim: Roles not found or not an array")
	}

	// Convert roles to []string
	var roles []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			uc.logger.Error("Invalid role type", zap.Any("role", role), zap.String("ip", ip))
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		roles = append(roles, roleStr)
	}

	// Fetch user details based on the user ID from the claims
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: &userID})
	if err != nil {
		uc.logger.Error("UserDatabase operation failed", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		uc.logger.Warn("User not found", zap.String("user_id", userID), zap.String("ip", ip))
		return nil, ErrUserNotFound
	}

	// Compare the refresh token with the stored hashed refresh token
	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		uc.logger.Error("Token comparison failed", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if !match {
		uc.logger.Warn("Refresh token mismatch", zap.String("user_id", userID), zap.String("ip", ip))
		return nil, ErrInvalidToken
	}

	// Create payload map for token generation
	payload := map[string]interface{}{
		"UserID": userID, // Ensure the correct field here
		"Roles":  roles,
	}

	// Generate new access and refresh tokens
	accessToken, accessExpiresAt, err := uc.jwtManager.GenerateAccessToken(payload)
	if err != nil {
		uc.logger.Error("Access token generation failed", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("%w: access token", ErrTokenGeneration)
	}

	newRefreshToken, refreshExpiresAt, err := uc.jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		uc.logger.Error("Refresh token generation failed", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("%w: refresh token", ErrTokenGeneration)
	}

	// Hash the generated access and refresh tokens
	hashedAccessToken, err := authutils.GenHash(uc.tokenPepper, accessToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash access token", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(uc.tokenPepper, newRefreshToken, nil)
	if err != nil {
		uc.logger.Error("Failed to hash refresh token", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("failed to secure tokens: %w", err)
	}

	// Update the user's stored tokens in the database
	_, err = uc.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		uc.logger.Error("Failed to update user tokens", zap.String("user_id", userID), zap.Error(err), zap.String("ip", ip))
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	uc.logger.Info("Tokens refreshed successfully", zap.String("user_id", userID), zap.Time("access_token_expiry", accessExpiresAt), zap.Time("refresh_token_expiry", refreshExpiresAt), zap.String("ip", ip))

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}
