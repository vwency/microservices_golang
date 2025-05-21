package logout

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value("ip").(string); ok {
		return ip
	}
	return "unknown"
}

func Logout(
	dbClient databasev1.DatabaseInitServiceClient,
	logger log.Logger,
	tokenPepper string,
	ctx context.Context,
	req *authv1.LogoutRequest,
) (*authv1.LogoutResponse, error) {

	ip := getIPFromContext(ctx)
	_ = level.Info(logger).Log(
		"msg", "Attempting logout",
		"username", req.Username,
		"ip", ip,
	)

	if req.Username == "" || req.AccessToken == "" {
		_ = level.Warn(logger).Log(
			"msg", "Missing logout credentials",
			"username", req.Username,
		)
		return nil, error_hndl.NewError(codes.InvalidArgument, "username and access_token are required")
	}

	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &req.Username,
	})
	if err != nil {
		_ = level.Error(logger).Log(
			"msg", "Failed to get user for logout",
			"username", req.Username,
			"err", err,
		)
		return nil, error_hndl.NewError(codes.NotFound, "failed to get user")
	}

	if !getUserResp.Found {
		_ = level.Warn(logger).Log(
			"msg", "User not found during logout attempt",
			"username", req.Username,
		)
		return nil, error_hndl.NewError(codes.NotFound, "user not found")
	}

	if isTokenEmptyOrNone(getUserResp.HashedAccessToken) {
		_ = level.Info(logger).Log(
			"msg", "User already logged out",
			"user_id", getUserResp.UserId,
			"username", req.Username,
		)
		return &authv1.LogoutResponse{
			Success: true,
			Message: "already logged out",
		}, nil
	}

	match, err := authutils.ComparePasswordAndHash(tokenPepper, req.AccessToken, getUserResp.HashedAccessToken)
	if err != nil {
		_ = level.Error(logger).Log(
			"msg", "Access token comparison failed",
			"username", req.Username,
			"err", err,
		)
		return nil, error_hndl.NewError(codes.Unauthenticated, "invalid access token")
	}
	if !match {
		_ = level.Warn(logger).Log(
			"msg", "Invalid access token provided",
			"username", req.Username,
		)
		return nil, error_hndl.NewError(codes.Unauthenticated, "access token mismatch")
	}

	// Функция для обновления токенов с двумя попытками:
	tryUpdateTokens := func(hashedRefresh, hashedAccess string) error {
		_, err := dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
			UserId:             getUserResp.UserId,
			HashedRefreshToken: hashedRefresh,
			HashedAccessToken:  hashedAccess,
		})
		return err
	}

	err = tryUpdateTokens("none", "none")
	if err != nil {
		_ = level.Error(logger).Log(
			"msg", "Failed to clear user tokens during logout",
			"user_id", getUserResp.UserId,
			"err", err,
		)

		if isTokenValidationError(err) {
			_ = level.Warn(logger).Log(
				"msg", "Retrying logout with empty tokens instead of 'none'",
				"user_id", getUserResp.UserId,
			)
			err = tryUpdateTokens("", "")
			if err != nil {
				_ = level.Error(logger).Log(
					"msg", "Fallback logout failed",
					"user_id", getUserResp.UserId,
					"err", err,
				)
				return nil, error_hndl.NewError(codes.NotFound, "logout failed")
			}
		} else {
			return nil, error_hndl.NewError(codes.NotFound, "logout failed")
		}
	}

	_ = level.Info(logger).Log(
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
