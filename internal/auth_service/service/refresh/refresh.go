package refresh

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

func getIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value("ip").(string); ok {
		return ip
	}
	return "unknown"
}

func Refresh(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger log.Logger,
	tokenPepper string,
	ctx context.Context,
	req *authv1.RefreshRequest,
) (*authv1.RefreshResponse, error) {

	ip := getIPFromContext(ctx)
	_ = level.Info(logger).Log("msg", "Attempting token refresh", "ip", ip)

	refreshToken := req.GetRefreshToken()
	claims, err := jwtManager.ValidateToken(refreshToken)
	if err != nil {
		_ = level.Warn(logger).Log("msg", "Invalid refresh token", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: %v", error_hndl.ErrInvalidToken, err)
	}

	userID, ok := claims["UserID"].(string)
	if !ok {
		_ = level.Error(logger).Log("msg", "Invalid claim: UserID missing", "ip", ip)
		return nil, errors.New("invalid claim: UserID not found or not a string")
	}

	rolesIface, ok := claims["Roles"].([]interface{})
	if !ok {
		_ = level.Error(logger).Log("msg", "Invalid claim: Roles missing or malformed", "ip", ip)
		return nil, errors.New("invalid claim: Roles not found or not an array")
	}

	var roles []string
	for _, role := range rolesIface {
		if r, ok := role.(string); ok {
			roles = append(roles, r)
		} else {
			_ = level.Error(logger).Log("msg", "Invalid role type", "role", role, "ip", ip)
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
	}

	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		_ = level.Error(logger).Log("msg", "Database error while fetching user", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if !getUserResp.Found {
		_ = level.Warn(logger).Log("msg", "User not found", "user_id", userID, "ip", ip)
		return nil, error_hndl.ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Refresh token comparison failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		_ = level.Warn(logger).Log("msg", "Refresh token mismatch", "user_id", userID, "ip", ip)
		return nil, error_hndl.ErrInvalidToken
	}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(payload)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Access token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: access token", error_hndl.ErrTokenGeneration)
	}

	newRefreshToken, _, err := jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Refresh token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: refresh token", error_hndl.ErrTokenGeneration)
	}

	hashedAccessToken, err := authutils.GenHash(tokenPepper, accessToken, nil)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Access token hashing failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(tokenPepper, newRefreshToken, nil)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Refresh token hashing failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	_, err = dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	})
	if err != nil {
		_ = level.Error(logger).Log("msg", "Failed to update tokens", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	_ = level.Info(logger).Log("msg", "Token refresh successful", "user_id", userID, "ip", ip)

	return &authv1.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
