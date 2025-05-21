package register

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
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

// Register реализует логику регистрации пользователя
func Register(
	dbClient databasev1.DatabaseInitServiceClient,
	logger interface {
		Log(keyvals ...interface{}) error
	},
	jwtManager interface {
		GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error)
		GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error)
	},
	tokenPepper string,
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {

	tracer := otel.Tracer("auth_service")
	ctx, span := tracer.Start(ctx, "RegisterService",
		trace.WithAttributes(
			otelAttr.String("username", req.GetUsername()),
			otelAttr.String("email", req.GetEmail()),
			otelAttr.Int("password.length", len(req.Password)),
		))
	defer span.End()

	level.Warn(logger).Log("Attempting registration",
		"username", req.Username, "ip", "TODO: getIPFromContext(ctx)")

	if err := validateRegisterInput(ctx, tracer, req); err != nil {
		return nil, err
	}

	userID := uuid.New().String()

	type tokenData struct {
		hashedPassword     string
		accessToken        string
		accessExpiresAt    time.Time
		refreshToken       string
		hashedAccessToken  string
		hashedRefreshToken string
	}

	data := &tokenData{}
	errChan := make(chan error, 5)
	var wg sync.WaitGroup

	ctx, parallelSpan := tracer.Start(ctx, "ParallelOperations")
	defer parallelSpan.End()

	wg.Add(5)

	// 1. Hash password
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "HashPassword")
		defer span.End()

		start := time.Now()
		var err error
		data.hashedPassword, err = authutils.GenHash(req.Username, req.Password, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			authErr := error_hndl.NewError(codes.Internal, "failed to hash password", err)
			errChan <- authErr
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			return
		}
		span.SetStatus(otelCodes.Ok, "password hashed successfully")
	}()

	// 2. Generate access token
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateAccessToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		var err error
		data.accessToken, data.accessExpiresAt, err = jwtManager.GenerateAccessToken(payload)
		span.SetAttributes(
			otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()),
			otelAttr.String("token.expires_at", data.accessExpiresAt.String()),
		)

		if err != nil {
			authErr := error_hndl.NewError(codes.Internal, "failed to generate access token", err)
			errChan <- authErr
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			return
		}
		span.SetStatus(otelCodes.Ok, "access token generated")
	}()

	// 3. Generate refresh token
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "GenerateRefreshToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		var err error
		data.refreshToken, _, err = jwtManager.GenerateRefreshToken(payload)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			authErr := error_hndl.NewError(codes.Internal, "failed to generate refresh token", err)
			errChan <- authErr
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			return
		}
		span.SetStatus(otelCodes.Ok, "refresh token generated")
	}()

	// 4. Hash access token (ждём генерации access token)
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "HashAccessToken")
		defer span.End()

		start := time.Now()

		// Ждем, пока accessToken не появится или ошибка не случится
		for {
			if data.accessToken != "" {
				break
			}
			select {
			case err := <-errChan:
				// Ошибка произошла в другой горутине — прекращаем
				span.RecordError(err)
				span.SetStatus(otelCodes.Error, "stopped due to prior error")
				return
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}

		var err error
		data.hashedAccessToken, err = authutils.GenHash(tokenPepper, data.accessToken, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			authErr := error_hndl.NewError(codes.Internal, "failed to hash access token", err)
			errChan <- authErr
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			return
		}
		span.SetStatus(otelCodes.Ok, "access token hashed")
	}()

	// 5. Hash refresh token (ждём генерации refresh token)
	go func() {
		defer wg.Done()
		_, span := tracer.Start(ctx, "HashRefreshToken")
		defer span.End()

		start := time.Now()

		// Ждем, пока refreshToken не появится или ошибка не случится
		for {
			if data.refreshToken != "" {
				break
			}
			select {
			case err := <-errChan:
				span.RecordError(err)
				span.SetStatus(otelCodes.Error, "stopped due to prior error")
				return
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}

		var err error
		data.hashedRefreshToken, err = authutils.GenHash(tokenPepper, data.refreshToken, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			authErr := error_hndl.NewError(codes.Internal, "failed to hash refresh token", err)
			errChan <- authErr
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			return
		}
		span.SetStatus(otelCodes.Ok, "refresh token hashed")
	}()

	wg.Wait()

	// Проверка ошибок из горутин
	select {
	case err := <-errChan:
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "registration failed during token operations")
		return nil, err
	default:
		close(errChan)
	}

	// Запрос на добавление пользователя в базу
	dbReq := &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     data.hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  data.hashedAccessToken,
		HashedRefreshToken: data.hashedRefreshToken,
		UserId:             &userID,
	}

	ctx, dbSpan := tracer.Start(ctx, "DatabaseAddUser")
	defer dbSpan.End()

	start := time.Now()
	_, err := dbClient.AddUser(ctx, dbReq)
	dbSpan.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		var authError error

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				authError = error_hndl.ErrUserAlreadyExists
			default:
				authError = error_hndl.NewError(st.Code(), "database operation failed", st.Message())
			}
		} else {
			authError = error_hndl.NewError(codes.Internal, "database operation failed", err)
		}

		dbSpan.RecordError(authError)
		dbSpan.SetStatus(otelCodes.Error, authError.Error())
		span.RecordError(authError)
		span.SetStatus(otelCodes.Error, "registration failed: database operation error")
		return nil, authError
	}

	logger.Log("level", "info", "msg", "User registered successfully",
		"user_id", userID, "username", req.Username, "email", req.Email)

	dbSpan.SetStatus(otelCodes.Ok, "user added to database")
	span.SetStatus(otelCodes.Ok, "registration successful")

	return &authv1.RegisterResponse{
		AccessToken:  data.accessToken,
		RefreshToken: data.refreshToken,
		ExpiresAt:    data.accessExpiresAt.Unix(),
	}, nil
}
