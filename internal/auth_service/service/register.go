package service

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.opentelemetry.io/otel"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateRegisterInput(ctx context.Context, tracer trace.Tracer, req *authv1.RegisterRequest) error {
	start := time.Now()
	_, span := tracer.Start(ctx, "InputValidation")
	defer span.End()

	if req.Username == "" || req.Password == "" || req.Email == "" {
		err := NewAuthError(codes.InvalidArgument, "username, password and email are required")
		span.SetAttributes(
			otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
			otelAttr.Bool("validation_passed", false),
		)
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
		otelAttr.Bool("validation_passed", true),
	)
	return nil
}

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "RegisterService",
		trace.WithAttributes(
			otelAttr.String("username", req.GetUsername()),
			otelAttr.String("email", req.GetEmail()),
			otelAttr.Int("password.length", len(req.Password)),
		))
	defer span.End()

	s.logger.Log("level", "info", "msg", "Attempting registration", "username", req.Username, "ip", getIPFromContext(ctx))

	if err := validateRegisterInput(ctx, tracer, req); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
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
		wg              sync.WaitGroup
	)

	ctx, parallelSpan := tracer.Start(ctx, "ParallelOperations")
	wg.Add(3)

	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "HashPassword",
			trace.WithAttributes(
				otelAttr.Int("input.length", len(req.Password)),
			))
		defer span.End()

		start := time.Now()
		hashedPassword, err = authutils.GenHash(req.Username, req.Password, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))
		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to hash password", err)
			span.RecordError(err)
		}
	}()

	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateAccessToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		accessToken, accessExpiresAt, err = s.jwtManager.GenerateAccessToken(payload)
		span.SetAttributes(
			otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()),
			otelAttr.String("token.expires_at", accessExpiresAt.String()),
		)
		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to generate access token", err)
			span.RecordError(err)
		}
	}()

	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateRefreshToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		refreshToken, _, err = s.jwtManager.GenerateRefreshToken(payload)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))
		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to generate refresh token", err)
			span.RecordError(err)
		}
	}()

	wg.Wait()
	parallelSpan.End()
	close(errChan)

	for e := range errChan {
		if e != nil {
			span.RecordError(e)
			span.SetStatus(otelCodes.Error, "parallel operations failed")
			return nil, e
		}
	}

	var hashedAccessToken, hashedRefreshToken string
	{
		_, hashSpan := tracer.Start(ctx, "HashTokens",
			trace.WithAttributes(
				otelAttr.Int("access_token.length", len(accessToken)),
				otelAttr.Int("refresh_token.length", len(refreshToken)),
			))
		defer hashSpan.End()

		start := time.Now()
		hashedAccessToken, err = authutils.GenHash(s.tokenPepper, accessToken, nil)
		hashSpan.SetAttributes(otelAttr.Int64("access_token.hash_duration.ns", time.Since(start).Nanoseconds()))
		if err != nil {
			err = NewAuthError(codes.Internal, "failed to hash access token", err)
			hashSpan.RecordError(err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			return nil, err
		}

		start = time.Now()
		hashedRefreshToken, err = authutils.GenHash(s.tokenPepper, refreshToken, nil)
		hashSpan.SetAttributes(otelAttr.Int64("refresh_token.hash_duration.ns", time.Since(start).Nanoseconds()))
		if err != nil {
			err = NewAuthError(codes.Internal, "failed to hash refresh token", err)
			hashSpan.RecordError(err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			return nil, err
		}
	}

	dbReq := &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
		UserId:             &userID,
	}

	ctx, dbSpan := tracer.Start(ctx, "DatabaseAddUser",
		trace.WithAttributes(
			otelAttr.String("db.operation", "AddUser"),
			otelAttr.Int("request.size.bytes", len(req.Username)+len(hashedPassword)+len(req.Email)+
				len(hashedAccessToken)+len(hashedRefreshToken)+len(userID)),
		))
	defer dbSpan.End()

	start := time.Now()
	_, err = s.dbClient.AddUser(ctx, dbReq)
	dbSpan.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				err = ErrUserAlreadyExists
			default:
				err = NewAuthError(st.Code(), "database operation failed", st.Message())
			}
		} else {
			err = NewAuthError(codes.Internal, "database operation failed", err)
		}
		dbSpan.RecordError(err)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	s.logger.Log("level", "info", "msg", "User registered successfully", "user_id", userID, "username", req.Username, "email", req.Email)
	span.SetStatus(otelCodes.Ok, "registration successful")

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
