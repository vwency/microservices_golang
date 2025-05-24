package login

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

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

var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrUserNotFound       = fmt.Errorf("user not found")
)

func Login(
	dbClient databasev1.DatabaseInitServiceClient,
	jwtManager *jwt.JWTManager,
	logger log.Logger,
	tokenPepper string,
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "Login")
	defer span.End()

	username := req.GetUsername()
	password := req.GetPassword()
	ip := getIPFromContext(ctx)

	span.SetAttributes(
		attribute.String("username", username),
		attribute.String("ip", ip),
	)

	logger.Log("message", "Attempting login", "username", username, "ip", ip)

	if username == "" || password == "" {
		err := ErrInvalidCredentials
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Log("message", "Empty credentials provided", "username", username)
		return nil, err
	}

	span.AddEvent("Fetching user from database")
	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: &username})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "DB error")
		logger.Log("message", "UserDatabase operation failed", "username", username, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		err := ErrUserNotFound
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		logger.Log("message", "User not found", "username", username)
		return nil, err
	}

	span.SetAttributes(attribute.String("user_id", getUserResp.UserId))

	span.AddEvent("Comparing password")
	match, err := authutils.ComparePasswordAndHash(username, password, getUserResp.HashedPassword)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Password comparison failed")
		logger.Log("message", "Password comparison failed", "username", username, "error", err)
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		err := ErrInvalidCredentials
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid password")
		logger.Log("message", "Invalid password", "username", username, "ip", ip)
		return nil, err
	}

	userID := getUserResp.UserId
	roles := []interface{}{"user"}

	span.AddEvent("Generating tokens")
	payload := map[string]interface{}{"UserID": userID, "Roles": roles}

	accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Access token generation failed")
		logger.Log("message", "Failed to generate access token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, _, err := jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Refresh token generation failed")
		logger.Log("message", "Failed to generate refresh token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	span.AddEvent("Hashing tokens")
	hashedAccessToken, err := authutils.GenHash(tokenPepper, accessToken, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Access token hashing failed")
		logger.Log("message", "Failed to hash access token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(tokenPepper, refreshToken, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Refresh token hashing failed")
		logger.Log("message", "Failed to hash refresh token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	span.AddEvent("Updating user tokens in DB")
	_, err = dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update tokens in DB")
		logger.Log("message", "Failed to update user tokens", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	span.SetStatus(codes.Ok, "Login successful")
	span.SetAttributes(attribute.String("user_id", userID), attribute.String("access_expires_at", accessExpiresAt.String()))
	logger.Log("message", "Login successful", "user_id", userID, "username", username, "token_expiry", accessExpiresAt)

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
