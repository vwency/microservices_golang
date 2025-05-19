package login

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
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

	username := req.GetUsername()
	password := req.GetPassword()

	logger.Log("message", "Attempting login for user",
		"username", username,
		"ip", getIPFromContext(ctx))

	if username == "" || password == "" {
		logger.Log("message", "Empty credentials provided", "username", username)
		return nil, ErrInvalidCredentials
	}

	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &username,
	})
	if err != nil {
		logger.Log("message", "UserDatabase operation failed", "username", username, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !getUserResp.Found {
		logger.Log("message", "User not found", "username", username)
		return nil, ErrUserNotFound
	}

	match, err := authutils.ComparePasswordAndHash(username, password, getUserResp.HashedPassword)
	if err != nil {
		logger.Log("message", "Password comparison failed", "username", username, "error", err)
		return nil, fmt.Errorf("authentication error: %w", err)
	}
	if !match {
		logger.Log("message", "Invalid password", "username", username, "ip", getIPFromContext(ctx))
		return nil, ErrInvalidCredentials
	}

	userID := getUserResp.UserId
	roles := []interface{}{"user"}

	payload := map[string]interface{}{
		"UserID": userID,
		"Roles":  roles,
	}

	accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(payload)
	if err != nil {
		logger.Log("message", "Failed to generate access token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, _, err := jwtManager.GenerateRefreshToken(payload)
	if err != nil {
		logger.Log("message", "Failed to generate refresh token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	hashedAccessToken, err := authutils.GenHash(tokenPepper, accessToken, nil)
	if err != nil {
		logger.Log("message", "Failed to hash access token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(tokenPepper, refreshToken, nil)
	if err != nil {
		logger.Log("message", "Failed to hash refresh token", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to hash refresh token: %v", err)
	}

	_, err = dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
		UserId:             userID,
		HashedRefreshToken: hashedRefreshToken,
		HashedAccessToken:  hashedAccessToken,
	})
	if err != nil {
		logger.Log("message", "Failed to update user tokens", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to update tokens: %w", err)
	}

	logger.Log("message", "Login successful", "user_id", userID, "username", username, "token_expiry", accessExpiresAt)

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
