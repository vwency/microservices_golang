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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
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
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "Logout")
	defer span.End()

	username := req.GetUsername()
	ip := getIPFromContext(ctx)

	span.SetAttributes(
		attribute.String("username", username),
		attribute.String("ip", ip),
	)

	_ = level.Info(logger).Log(
		"msg", "Attempting logout",
		"username", username,
		"ip", ip,
	)

	if username == "" || req.AccessToken == "" {
		err := error_hndl.ErrInvalidUsernameFormat
		span.RecordError(err)
		span.SetStatus(codes.Error, "Missing credentials")
		_ = level.Warn(logger).Log(
			"msg", "Missing logout credentials",
			"username", username,
		)
		return nil, err
	}

	span.AddEvent("Fetching user from database")
	getUserResp, err := dbClient.GetUser(ctx, &databasev1.GetUserRequest{
		Username: &username,
	})
	if err != nil {
		span.RecordError(err)
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				span.SetStatus(codes.Error, "User not found")
				_ = level.Warn(logger).Log(
					"msg", "User not found during logout attempt (from dbClient)",
					"username", username,
					"err", err,
				)
				return nil, error_hndl.ErrUserNotFound
			case grpcCodes.Unavailable, grpcCodes.DeadlineExceeded:
				span.SetStatus(codes.Error, "Database unavailable")
				_ = level.Error(logger).Log(
					"msg", "Database unavailable or timeout during logout",
					"username", username,
					"err", err,
				)
				return nil, error_hndl.ErrDatabaseFailure
			default:
				span.SetStatus(codes.Error, "Database error")
				_ = level.Error(logger).Log(
					"msg", "Failed to get user for logout",
					"username", username,
					"err", err,
				)
				return nil, error_hndl.ErrDatabaseFailure
			}
		}
		span.SetStatus(codes.Error, "Unknown database error")
		_ = level.Error(logger).Log(
			"msg", "Unknown error during GetUser",
			"username", username,
			"err", err,
		)
		return nil, error_hndl.ErrDatabaseFailure
	}

	if !getUserResp.Found {
		err := error_hndl.ErrUserNotFound
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		_ = level.Warn(logger).Log(
			"msg", "User not found during logout attempt (response found=false)",
			"username", username,
		)
		return nil, err
	}

	span.SetAttributes(attribute.String("user_id", getUserResp.UserId))

	if isTokenEmptyOrNone(getUserResp.HashedAccessToken) {
		span.AddEvent("User already logged out")
		_ = level.Info(logger).Log(
			"msg", "User already logged out",
			"user_id", getUserResp.UserId,
			"username", username,
		)
		return &authv1.LogoutResponse{
			Success: true,
			Message: "already logged out",
		}, nil
	}

	span.AddEvent("Validating access token")
	match, err := authutils.ComparePasswordAndHash(tokenPepper, req.AccessToken, getUserResp.HashedAccessToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token validation failed")
		_ = level.Error(logger).Log(
			"msg", "Access token comparison failed",
			"username", username,
			"err", err,
		)
		return nil, error_hndl.ErrInvalidToken
	}
	if !match {
		err := error_hndl.ErrAccessTokenMismatch
		span.RecordError(err)
		span.SetStatus(codes.Error, "Token mismatch")
		_ = level.Warn(logger).Log(
			"msg", "Invalid access token provided",
			"username", username,
		)
		return nil, err
	}

	tryUpdateTokens := func(hashedRefresh, hashedAccess string) error {
		_, err := dbClient.UpdateUser(ctx, &databasev1.UpdateUserRequest{
			UserId:             getUserResp.UserId,
			HashedRefreshToken: hashedRefresh,
			HashedAccessToken:  hashedAccess,
		})
		return err
	}

	span.AddEvent("Clearing user tokens")
	err = tryUpdateTokens("none", "none")
	if err != nil {
		span.RecordError(err)
		_ = level.Error(logger).Log(
			"msg", "Failed to clear user tokens during logout",
			"user_id", getUserResp.UserId,
			"err", err,
		)

		if isTokenRelatedError(err) {
			span.AddEvent("Retrying with empty tokens")
			_ = level.Warn(logger).Log(
				"msg", "Retrying logout with empty tokens instead of 'none'",
				"user_id", getUserResp.UserId,
			)
			err = tryUpdateTokens("", "")
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "Logout failed")
				_ = level.Error(logger).Log(
					"msg", "Fallback logout failed",
					"user_id", getUserResp.UserId,
					"err", err,
				)
				return nil, error_hndl.ErrLogoutFailed
			}
		} else {
			span.SetStatus(codes.Error, "Logout failed")
			return nil, error_hndl.ErrLogoutFailed
		}
	}

	span.SetStatus(codes.Ok, "Logout successful")
	_ = level.Info(logger).Log(
		"msg", "User logged out successfully",
		"user_id", getUserResp.UserId,
		"username", username,
	)

	return &authv1.LogoutResponse{
		Success: true,
		Message: "logged out successfully",
	}, nil
}

func isTokenRelatedError(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	if st.Code() == grpcCodes.InvalidArgument {
		msg := strings.ToLower(st.Message())
		if strings.Contains(msg, "token") {
			return true
		}
	}
	return false
}

func isTokenEmptyOrNone(token string) bool {
	return token == "" || strings.EqualFold(token, "none")
}
