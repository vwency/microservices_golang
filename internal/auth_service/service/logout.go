package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log/level"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	ip := getIPFromContext(ctx)
	_ = level.Info(s.logger).Log(
		"msg", "Attempting logout",
		"username", req.Username,
		"ip", ip,
	)

	if req.Username == "" || req.AccessToken == "" {
		_ = level.Warn(s.logger).Log(
			"msg", "Missing logout credentials",
			"username", req.Username,
		)
		return nil, status.Error(codes.InvalidArgument, "username and access_token are required")
	}

	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &req.Username,
	})
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to get user for logout",
			"username", req.Username,
			"err", err,
		)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if !getUserResp.Found {
		_ = level.Warn(s.logger).Log(
			"msg", "User not found during logout attempt",
			"username", req.Username,
		)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if isTokenEmptyOrNone(getUserResp.HashedAccessToken) {
		_ = level.Info(s.logger).Log(
			"msg", "User already logged out",
			"user_id", getUserResp.UserId,
			"username", req.Username,
		)
		return &authv1.LogoutResponse{
			Success: true,
			Message: "already logged out",
		}, nil
	}

	match, err := authutils.ComparePasswordAndHash(s.tokenPepper, req.AccessToken, getUserResp.HashedAccessToken)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Access token comparison failed",
			"username", req.Username,
			"err", err,
		)
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
	}

	if !match {
		_ = level.Warn(s.logger).Log(
			"msg", "Invalid access token provided",
			"username", req.Username,
		)
		return nil, status.Error(codes.Unauthenticated, "access token mismatch")
	}

	logoutRequest := &databasev1.UpdateUserRequest{
		UserId:             getUserResp.UserId,
		HashedRefreshToken: "none",
		HashedAccessToken:  "none",
	}

	_, err = s.dbClient.UpdateUser(ctx, logoutRequest)
	if err != nil {
		_ = level.Error(s.logger).Log(
			"msg", "Failed to clear user tokens during logout",
			"user_id", getUserResp.UserId,
			"err", err,
		)

		if isTokenValidationError(err) {
			_ = level.Warn(s.logger).Log(
				"msg", "Retrying logout with empty tokens instead of 'none'",
				"user_id", getUserResp.UserId,
			)

			_, err = s.dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
				UserId:             getUserResp.UserId,
				HashedRefreshToken: "",
				HashedAccessToken:  "",
			})
			if err != nil {
				_ = level.Error(s.logger).Log(
					"msg", "Fallback logout failed",
					"user_id", getUserResp.UserId,
					"err", err,
				)
				return nil, fmt.Errorf("fallback logout failed: %v", err)
			}
		} else {
			return nil, fmt.Errorf("logout failed: %v", err)
		}
	}

	_ = level.Info(s.logger).Log(
		"msg", "User logged out successfully",
		"user_id", getUserResp.UserId,
		"username", req.Username,
	)

	return &authv1.LogoutResponse{
		Success: true,
		Message: "logged out successfully",
	}, nil
}

func isTokenValidationError(err error) bool {
	if status, ok := status.FromError(err); ok {
		return status.Code() == codes.InvalidArgument &&
			(strings.Contains(status.Message(), "token") ||
				strings.Contains(status.Message(), "Token"))
	}
	return false
}

func isTokenEmptyOrNone(token string) bool {
	return token == "" || strings.EqualFold(token, "none")
}
