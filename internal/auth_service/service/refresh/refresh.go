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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "Refresh")
	defer span.End()

	ip := getIPFromContext(ctx)
	span.SetAttributes(
		attribute.String("ip", ip),
	)

	_ = level.Info(logger).Log("msg", "Attempting token refresh", "ip", ip)

	refreshToken := req.GetRefreshToken()
	claims, err := jwtManager.ValidateToken(refreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid refresh token")
		_ = level.Warn(logger).Log("msg", "Invalid refresh token", "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: %v", error_hndl.ErrInvalidToken, err)
	}

	userID, ok := claims["UserID"].(string)
	if !ok {
		err := errors.New("invalid claim: UserID not found or not a string")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid claim: UserID")
		_ = level.Error(logger).Log("msg", "Invalid claim: UserID missing", "ip", ip)
		return nil, err
	}

	span.SetAttributes(attribute.String("user_id", userID))

	rolesIface, ok := claims["Roles"].([]interface{})
	if !ok {
		err := errors.New("invalid claim: Roles not found or not an array")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid claim: Roles")
		_ = level.Error(logger).Log("msg", "Invalid claim: Roles missing or malformed", "ip", ip)
		return nil, err
	}

	var roles []string
	for _, role := range rolesIface {
		if r, ok := role.(string); ok {
			roles = append(roles, r)
		} else {
			err := fmt.Errorf("invalid role type: %v", role)
			span.RecordError(err)
			span.SetStatus(codes.Error, "Invalid role type")
			_ = level.Error(logger).Log("msg", "Invalid role type", "role", role, "ip", ip)
			return nil, err
		}
	}

	span.AddEvent("Fetching user from database")
	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Database error")
		_ = level.Error(logger).Log("msg", "Database error while fetching user", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if !getUserResp.Found {
		err := error_hndl.ErrUserNotFound
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		_ = level.Warn(logger).Log("msg", "User not found", "user_id", userID, "ip", ip)
		return nil, err
	}

	span.AddEvent("Validating refresh token")
	match, err := authutils.ComparePasswordAndHash(tokenPepper, refreshToken, getUserResp.HashedRefreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token validation failed")
		_ = level.Error(logger).Log("msg", "Refresh token comparison failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		err := error_hndl.ErrInvalidToken
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token mismatch")
		_ = level.Warn(logger).Log("msg", "Refresh token mismatch", "user_id", userID, "ip", ip)
		return nil, err
	}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	span.AddEvent("Generating new tokens")
	accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Access token generation failed")
		_ = level.Error(logger).Log("msg", "Access token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: access token", error_hndl.ErrTokenGeneration)
	}

	newRefreshToken, _, err := jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Refresh token generation failed")
		_ = level.Error(logger).Log("msg", "Refresh token generation failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("%w: refresh token", error_hndl.ErrTokenGeneration)
	}

	span.AddEvent("Hashing tokens")
	hashedAccessToken, err := authutils.GenHash(tokenPepper, accessToken, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Access token hashing failed")
		_ = level.Error(logger).Log("msg", "Access token hashing failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	hashedRefreshToken, err := authutils.GenHash(tokenPepper, newRefreshToken, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Refresh token hashing failed")
		_ = level.Error(logger).Log("msg", "Refresh token hashing failed", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	span.AddEvent("Updating tokens in database")
	_, err = dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update tokens")
		_ = level.Error(logger).Log("msg", "Failed to update tokens", "user_id", userID, "err", err, "ip", ip)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	span.SetStatus(codes.Ok, "Token refresh successful")
	_ = level.Info(logger).Log("msg", "Token refresh successful", "user_id", userID, "ip", ip)

	return &authv1.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
