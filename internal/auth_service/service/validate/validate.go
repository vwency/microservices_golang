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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "ValidateAccessToken")
	defer span.End()

	token := req.GetAccessToken()
	if token == "" {
		err := errors.New("access token is required")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Missing access token")
		deps.Logger.Log("error", "Access token is required")
		return nil, err
	}

	span.SetAttributes(
		attribute.String("token", token[:3]+"***"), // Логируем только первые 3 символа для безопасности
	)

	span.AddEvent("Validating token structure")
	claims, err := deps.JwtManager.ValidateToken(token)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token validation failed")
		deps.Logger.Log("error", fmt.Sprintf("Failed to validate token: %v", err), "token", token)
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	userID, ok := claims["UserID"].(string)
	if !ok || userID == "" {
		err := errors.New("invalid token: userID missing or not a string")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid UserID in claims")
		deps.Logger.Log("error", "Invalid or missing UserID in claims", "claims", claims)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user_id", userID),
	)

	rolesRaw, ok := claims["Roles"].([]interface{})
	if !ok {
		err := errors.New("invalid token: roles missing or not an array")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid Roles in claims")
		deps.Logger.Log("error", "Invalid or missing Roles in claims", "claims", claims)
		return nil, err
	}

	var roles []string
	for _, r := range rolesRaw {
		roleStr, ok := r.(string)
		if !ok {
			err := fmt.Errorf("invalid role type: %v", r)
			span.RecordError(err)
			span.SetStatus(codes.Error, "Invalid role type")
			deps.Logger.Log("error", fmt.Sprintf("Invalid role type: %v", r))
			return nil, err
		}
		roles = append(roles, roleStr)
	}

	span.SetAttributes(
		attribute.StringSlice("roles", roles),
	)

	span.AddEvent("Fetching user from database")
	getUserResp, err := deps.DbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Database error")
		deps.Logger.Log("error", fmt.Sprintf("Failed to fetch user: %v", err), "userID", userID)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !getUserResp.Found {
		err := errors.New("user not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		deps.Logger.Log("warn", "User not found", "userID", userID)
		return nil, err
	}

	span.AddEvent("Comparing token hashes")
	match, err := authutils.ComparePasswordAndHash(deps.TokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token hash comparison failed")
		deps.Logger.Log("error", fmt.Sprintf("Token hash comparison failed: %v", err), "userID", userID)
		return nil, fmt.Errorf("token hash comparison failed: %w", err)
	}

	if !match {
		err := errors.New("invalid token: hash mismatch")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token hash mismatch")
		deps.Logger.Log("warn", "Token hash mismatch", "userID", userID)
		return nil, err
	}

	expiryFloat, ok := claims["exp"].(float64)
	if !ok {
		err := errors.New("invalid token: expiry missing or not a number")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid expiry in claims")
		deps.Logger.Log("error", "Invalid or missing expiry in claims", "claims", claims)
		return nil, err
	}
	expiry := int64(expiryFloat)

	span.SetAttributes(
		attribute.Int64("expiry", expiry),
	)

	if time.Now().Unix() > expiry {
		err := errors.New("access token has expired")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token expired")
		deps.Logger.Log("warn", "Token expired", "userID", userID, "expiry", expiry)
		return nil, err
	}

	span.SetStatus(codes.Ok, "Token validation successful")
	return &authv1.ValidateResponse{
		Valid:     true,
		UserId:    userID,
		Roles:     roles,
		ExpiresAt: expiry,
	}, nil
}
