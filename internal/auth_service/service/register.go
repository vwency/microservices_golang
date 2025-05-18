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

// validateRegisterInput validates registration input parameters
func validateRegisterInput(ctx context.Context, tracer trace.Tracer, req *authv1.RegisterRequest) error {
	ctx, span := tracer.Start(ctx, "InputValidation")
	defer span.End()

	start := time.Now()

	if req.Username == "" || req.Password == "" || req.Email == "" {
		err := NewAuthError(codes.InvalidArgument, "username, password and email are required")
		span.SetAttributes(
			otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
			otelAttr.Bool("validation_passed", false),
		)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return err
	}

	span.SetAttributes(
		otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
		otelAttr.Bool("validation_passed", true),
	)
	span.SetStatus(otelCodes.Ok, "validation passed")
	return nil
}

// Register handles user registration
func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "RegisterService",
		trace.WithAttributes(
			otelAttr.String("username", req.GetUsername()),
			otelAttr.String("email", req.GetEmail()),
			otelAttr.Int("password.length", len(req.Password)),
		))
	defer span.End()

	s.logger.Log("level", "info", "msg", "Attempting registration",
		"username", req.Username, "ip", getIPFromContext(ctx))

	// Validate input parameters
	if err := validateRegisterInput(ctx, tracer, req); err != nil {
		return nil, err // Error already logged and status set in validateRegisterInput
	}

	// Generate user ID
	userID := uuid.New().String()

	// Struct to hold all generated data
	type tokenData struct {
		hashedPassword     string
		accessToken        string
		accessExpiresAt    time.Time
		refreshToken       string
		hashedAccessToken  string
		hashedRefreshToken string
	}

	// Create channel to collect errors
	errChan := make(chan error, 5)
	data := &tokenData{}
	var wg sync.WaitGroup

	// Start parallel operations with parent span
	ctx, parallelSpan := tracer.Start(ctx, "ParallelOperations")
	defer parallelSpan.End()

	// Use 5 goroutines for full parallelization
	wg.Add(5)

	// 1. Hash password
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "HashPassword",
			trace.WithAttributes(otelAttr.Int("input.length", len(req.Password))))
		defer span.End()

		start := time.Now()
		var err error
		data.hashedPassword, err = authutils.GenHash(req.Username, req.Password, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))
		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to hash password", err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
		} else {
			span.SetStatus(otelCodes.Ok, "password hashed successfully")
		}
	}()

	// 2. Generate access token
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateAccessToken")
		defer span.End()

		start := time.Now()
		var err error
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		data.accessToken, data.accessExpiresAt, err = s.jwtManager.GenerateAccessToken(payload)
		span.SetAttributes(
			otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()),
			otelAttr.String("token.expires_at", data.accessExpiresAt.String()),
		)

		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to generate access token", err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
		} else {
			span.SetStatus(otelCodes.Ok, "access token generated")
		}
	}()

	// 3. Generate refresh token
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateRefreshToken")
		defer span.End()

		start := time.Now()
		var err error
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		data.refreshToken, _, err = s.jwtManager.GenerateRefreshToken(payload)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to generate refresh token", err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
		} else {
			span.SetStatus(otelCodes.Ok, "refresh token generated")
		}
	}()

	// 4. Hash access token - Now part of parallel operations
	go func() {
		defer wg.Done()

		// Wait for access token to be generated
		for data.accessToken == "" {
			time.Sleep(1 * time.Millisecond)
			// If another goroutine has already failed, exit early
			select {
			case <-errChan:
				return
			default:
			}
		}

		_, span := tracer.Start(ctx, "HashAccessToken",
			trace.WithAttributes(otelAttr.Int("token.length", len(data.accessToken))))
		defer span.End()

		start := time.Now()
		var err error
		data.hashedAccessToken, err = authutils.GenHash(s.tokenPepper, data.accessToken, nil)

		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to hash access token", err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
		} else {
			span.SetStatus(otelCodes.Ok, "access token hashed")
		}
	}()

	// 5. Hash refresh token - Now part of parallel operations
	go func() {
		defer wg.Done()

		// Wait for refresh token to be generated
		for data.refreshToken == "" {
			time.Sleep(1 * time.Millisecond)
			// If another goroutine has already failed, exit early
			select {
			case <-errChan:
				return
			default:
			}
		}

		_, span := tracer.Start(ctx, "HashRefreshToken",
			trace.WithAttributes(otelAttr.Int("token.length", len(data.refreshToken))))
		defer span.End()

		start := time.Now()
		var err error
		data.hashedRefreshToken, err = authutils.GenHash(s.tokenPepper, data.refreshToken, nil)

		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			errChan <- NewAuthError(codes.Internal, "failed to hash refresh token", err)
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
		} else {
			span.SetStatus(otelCodes.Ok, "refresh token hashed")
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Process any errors from the goroutines
	select {
	case err := <-errChan:
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "registration failed during token operations")
		return nil, err
	default:
		// No errors occurred, close the channel
		close(errChan)
	}

	// Create database request
	dbReq := &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     data.hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  data.hashedAccessToken,
		HashedRefreshToken: data.hashedRefreshToken,
		UserId:             &userID,
	}

	// Insert user into database
	ctx, dbSpan := tracer.Start(ctx, "DatabaseAddUser",
		trace.WithAttributes(
			otelAttr.String("db.operation", "AddUser"),
			otelAttr.String("user.id", userID),
			otelAttr.Int("request.size.bytes", len(req.Username)+len(data.hashedPassword)+len(req.Email)+
				len(data.hashedAccessToken)+len(data.hashedRefreshToken)+len(userID)),
		))
	defer dbSpan.End()

	start := time.Now()
	_, err := s.dbClient.AddUser(ctx, dbReq)
	dbSpan.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		var authError error

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				authError = ErrUserAlreadyExists
			default:
				authError = NewAuthError(st.Code(), "database operation failed", st.Message())
			}
		} else {
			authError = NewAuthError(codes.Internal, "database operation failed", err)
		}

		dbSpan.RecordError(authError)
		dbSpan.SetStatus(otelCodes.Error, authError.Error())
		span.RecordError(authError)
		span.SetStatus(otelCodes.Error, "registration failed: database operation error")
		return nil, authError
	}

	// Log successful registration
	s.logger.Log("level", "info", "msg", "User registered successfully",
		"user_id", userID, "username", req.Username, "email", req.Email)

	dbSpan.SetStatus(otelCodes.Ok, "user added to database")
	span.SetStatus(otelCodes.Ok, "registration successful")

	// Return response
	return &authv1.RegisterResponse{
		AccessToken:  data.accessToken,
		RefreshToken: data.refreshToken,
		ExpiresAt:    data.accessExpiresAt.Unix(),
	}, nil
}
