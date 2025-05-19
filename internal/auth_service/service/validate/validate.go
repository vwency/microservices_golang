package validate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

type ValidateDeps struct {
	DbClient    databasev1.DatabaseInitServiceClient
	JwtManager  JWTManagerInterface
	Logger      log.Logger
	TokenPepper string
}

type JWTManagerInterface interface {
	ValidateToken(token string) (map[string]interface{}, error)
}

func ValidateAccessToken(deps ValidateDeps, ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	token := req.GetAccessToken()
	if token == "" {
		return nil, errors.New("access token is required")
	}

	claims, err := deps.JwtManager.ValidateToken(token)
	if err != nil {
		deps.Logger.Log("error", fmt.Sprintf("Failed to validate token: %v", err), "token", token)
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	userID, ok := claims["UserID"].(string)
	if !ok || userID == "" {
		deps.Logger.Log("error", "Invalid or missing UserID in claims", "claims", claims)
		return nil, errors.New("invalid token: userID missing or not a string")
	}

	rolesRaw, ok := claims["Roles"].([]interface{})
	if !ok {
		deps.Logger.Log("error", "Invalid or missing Roles in claims", "claims", claims)
		return nil, errors.New("invalid token: roles missing or not an array")
	}

	var roles []string
	for _, r := range rolesRaw {
		roleStr, ok := r.(string)
		if !ok {
			deps.Logger.Log("error", fmt.Sprintf("Invalid role type: %v", r))
			return nil, fmt.Errorf("invalid role type: %v", r)
		}
		roles = append(roles, roleStr)
	}

	getUserResp, err := deps.DbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		deps.Logger.Log("error", fmt.Sprintf("Failed to fetch user: %v", err), "userID", userID)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !getUserResp.Found {
		deps.Logger.Log("warn", "User not found", "userID", userID)
		return nil, errors.New("user not found")
	}

	match, err := authutils.ComparePasswordAndHash(deps.TokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		deps.Logger.Log("error", fmt.Sprintf("Token hash comparison failed: %v", err), "userID", userID)
		return nil, fmt.Errorf("token hash comparison failed: %w", err)
	}

	if !match {
		deps.Logger.Log("warn", "Token hash mismatch", "userID", userID)
		return nil, errors.New("invalid token: hash mismatch")
	}

	expiryFloat, ok := claims["exp"].(float64)
	if !ok {
		deps.Logger.Log("error", "Invalid or missing expiry in claims", "claims", claims)
		return nil, errors.New("invalid token: expiry missing or not a number")
	}
	expiry := int64(expiryFloat)

	if time.Now().Unix() > expiry {
		deps.Logger.Log("warn", "Token expired", "userID", userID, "expiry", expiry)
		return nil, errors.New("access token has expired")
	}

	return &authv1.ValidateResponse{
		Valid:     true,
		UserId:    userID,
		Roles:     roles,
		ExpiresAt: expiry,
	}, nil
}
