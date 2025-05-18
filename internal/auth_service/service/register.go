package service

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.opentelemetry.io/otel"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "RegisterService",
		trace.WithAttributes(
			otelAttr.String("username", req.GetUsername()),
			otelAttr.String("email", req.GetEmail()),
		))
	defer span.End()

	level.Info(s.logger).Log(
		"msg", "Attempting registration",
		"username", req.Username,
		"ip", getIPFromContext(ctx),
	)

	if req.Username == "" || req.Password == "" || req.Email == "" {
		err := status.Error(grpcCodes.InvalidArgument, "username, password and email are required")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "validation failed")
		return nil, err
	}

	var (
		userID          = uuid.New().String()
		hashedPassword  string
		accessToken     string
		accessExpiresAt time.Time
		refreshToken    string
		err             error
		errChan         = make(chan error, 3)
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		_, span := tracer.Start(context.Background(), "HashPassword")
		defer span.End()
		hashedPassword, err = authutils.GenHash(req.Username, req.Password, nil)
		if err != nil {
			errChan <- err
			span.RecordError(err)
		}
	}()

	go func() {
		defer wg.Done()
		_, span := tracer.Start(context.Background(), "GenerateAccessToken")
		defer span.End()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		accessToken, accessExpiresAt, err = s.jwtManager.GenerateAccessToken(payload)
		if err != nil {
			errChan <- err
			span.RecordError(err)
		}
	}()

	go func() {
		defer wg.Done()
		_, span := tracer.Start(context.Background(), "GenerateRefreshToken")
		defer span.End()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		refreshToken, _, err = s.jwtManager.GenerateRefreshToken(payload)
		if err != nil {
			errChan <- err
			span.RecordError(err)
		}
	}()

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			span.RecordError(e)
			span.SetStatus(otelCodes.Error, "parallel operations failed")
			return nil, status.Errorf(grpcCodes.Internal, "operation failed: %v", e)
		}
	}

	ctx, hashSpan := tracer.Start(ctx, "HashTokens")
	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		hashSpan.RecordError(err)
		hashSpan.End()
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to hash tokens")
		return nil, status.Errorf(grpcCodes.Internal, "failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		hashSpan.RecordError(err)
		hashSpan.End()
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "failed to hash tokens")
		return nil, status.Errorf(grpcCodes.Internal, "failed to hash refresh token: %v", err)
	}
	hashSpan.End()

	ctx, dbSpan := tracer.Start(ctx, "DatabaseAddUser")
	defer dbSpan.End()
	if _, err := s.dbClient.AddUser(ctx, &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
		UserId:             &userID,
	}); err != nil {
		dbSpan.RecordError(err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "database operation failed")
		return nil, status.Errorf(grpcCodes.Internal, "failed to add user: %v", err)
	}

	level.Info(s.logger).Log(
		"msg", "User registered successfully",
		"user_id", userID,
		"username", req.Username,
	)

	span.SetStatus(otelCodes.Ok, "registration successful")
	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
