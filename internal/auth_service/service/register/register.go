package register

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
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

func (td *tokenData) reset() {
	td.hashedPassword = ""
	td.accessToken = ""
	td.refreshToken = ""
	td.hashedAccessToken = ""
	td.hashedRefreshToken = ""
	td.accessExpiresAt = time.Time{}
}

const (
	taskPassword = iota
	taskAccessToken
	taskRefreshToken
	taskHashAccess
	taskHashRefresh
)

// Оптимизированная функция с исправленной логикой
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

	// Получаем структуру данных из пула
	data := tokenDataPool.Get().(*tokenData)
	defer func() {
		data.reset()
		tokenDataPool.Put(data)
	}()

	// Получаем WaitGroup из пула
	wg := wgPool.Get().(*sync.WaitGroup)
	defer wgPool.Put(wg)

	ctx, parallelSpan := tracer.Start(ctx, "ParallelOperations")
	defer parallelSpan.End()

	// Канал для результатов с буферизацией
	resultChan := make(chan taskResult, 5)

	// Атомарный счетчик для отслеживания ошибок
	var hasError int32

	// Каналы синхронизации для зависимых задач
	accessTokenReady := make(chan struct{})
	refreshTokenReady := make(chan struct{})

	// Ограничиваем количество горутин
	maxGoroutines := runtime.NumCPU()
	if maxGoroutines > 3 {
		maxGoroutines = 3
	}
	semaphore := make(chan struct{}, maxGoroutines)

	// Хелпер для проверки контекста и ошибок
	checkError := func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			return atomic.LoadInt32(&hasError) != 0
		}
	}

	wg.Add(5)

	// 1. Хэширование пароля (критически важно - используем стандартные параметры)
	go func() {
		defer func() {
			wg.Done()
			<-semaphore
		}()
		semaphore <- struct{}{}

		if checkError() {
			resultChan <- taskResult{taskPassword, ctx.Err()}
			return
		}

		_, span := tracer.Start(ctx, "HashPassword")
		defer span.End()

		start := time.Now()
		hashedPassword, err := authutils.GenHash(req.Username, req.Password, nil)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			atomic.StoreInt32(&hasError, 1)
			authErr := error_hndl.NewError(codes.Internal, "failed to hash password", err)
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			resultChan <- taskResult{taskPassword, authErr}
			return
		}

		data.hashedPassword = hashedPassword
		span.SetStatus(otelCodes.Ok, "password hashed successfully")
		resultChan <- taskResult{taskPassword, nil}
	}()

	// 2. Генерация access token
	go func() {
		defer func() {
			wg.Done()
			<-semaphore
		}()
		semaphore <- struct{}{}

		if checkError() {
			resultChan <- taskResult{taskAccessToken, ctx.Err()}
			return
		}

		_, span := tracer.Start(ctx, "GenerateAccessToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(payload)
		span.SetAttributes(
			otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()),
		)

		if err != nil {
			atomic.StoreInt32(&hasError, 1)
			authErr := error_hndl.NewError(codes.Internal, "failed to generate access token", err)
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			resultChan <- taskResult{taskAccessToken, authErr}
			return
		}

		data.accessToken = accessToken
		data.accessExpiresAt = accessExpiresAt
		close(accessTokenReady)
		span.SetStatus(otelCodes.Ok, "access token generated")
		resultChan <- taskResult{taskAccessToken, nil}
	}()

	// 3. Генерация refresh token
	go func() {
		defer func() {
			wg.Done()
			<-semaphore
		}()
		semaphore <- struct{}{}

		if checkError() {
			resultChan <- taskResult{taskRefreshToken, ctx.Err()}
			return
		}

		_, span := tracer.Start(ctx, "GenerateRefreshToken")
		defer span.End()

		start := time.Now()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}

		refreshToken, _, err := jwtManager.GenerateRefreshToken(payload)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			atomic.StoreInt32(&hasError, 1)
			authErr := error_hndl.NewError(codes.Internal, "failed to generate refresh token", err)
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			resultChan <- taskResult{taskRefreshToken, authErr}
			return
		}

		data.refreshToken = refreshToken
		close(refreshTokenReady)
		span.SetStatus(otelCodes.Ok, "refresh token generated")
		resultChan <- taskResult{taskRefreshToken, nil}
	}()

	// 4. Хэширование access token (зависит от генерации токена)
	go func() {
		defer func() {
			wg.Done()
			<-semaphore
		}()
		semaphore <- struct{}{}

		// Ждем готовности access token
		select {
		case <-ctx.Done():
			resultChan <- taskResult{taskHashAccess, ctx.Err()}
			return
		case <-accessTokenReady:
		}

		if checkError() {
			resultChan <- taskResult{taskHashAccess, ctx.Err()}
			return
		}

		_, span := tracer.Start(ctx, "HashAccessToken")
		defer span.End()

		start := time.Now()
		hashedAccessToken, err := authutils.GenHash(tokenPepper, data.accessToken, ultraFastHashParams)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			atomic.StoreInt32(&hasError, 1)
			authErr := error_hndl.NewError(codes.Internal, "failed to hash access token", err)
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			resultChan <- taskResult{taskHashAccess, authErr}
			return
		}

		data.hashedAccessToken = hashedAccessToken
		span.SetStatus(otelCodes.Ok, "access token hashed")
		resultChan <- taskResult{taskHashAccess, nil}
	}()

	// 5. Хэширование refresh token (зависит от генерации токена)
	go func() {
		defer func() {
			wg.Done()
			<-semaphore
		}()
		semaphore <- struct{}{}

		// Ждем готовности refresh token
		select {
		case <-ctx.Done():
			resultChan <- taskResult{taskHashRefresh, ctx.Err()}
			return
		case <-refreshTokenReady:
		}

		if checkError() {
			resultChan <- taskResult{taskHashRefresh, ctx.Err()}
			return
		}

		_, span := tracer.Start(ctx, "HashRefreshToken")
		defer span.End()

		start := time.Now()
		hashedRefreshToken, err := authutils.GenHash(tokenPepper, data.refreshToken, ultraFastHashParams)
		span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

		if err != nil {
			atomic.StoreInt32(&hasError, 1)
			authErr := error_hndl.NewError(codes.Internal, "failed to hash refresh token", err)
			span.RecordError(authErr)
			span.SetStatus(otelCodes.Error, authErr.Error())
			resultChan <- taskResult{taskHashRefresh, authErr}
			return
		}

		data.hashedRefreshToken = hashedRefreshToken
		span.SetStatus(otelCodes.Ok, "refresh token hashed")
		resultChan <- taskResult{taskHashRefresh, nil}
	}()

	// Ждем завершения всех задач и обрабатываем результаты
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Обрабатываем результаты
	completedTasks := 0
	for result := range resultChan {
		if result.err != nil {
			span.RecordError(result.err)
			span.SetStatus(otelCodes.Error, "registration failed during parallel operations")
			return nil, result.err
		}
		completedTasks++
	}

	// Проверяем, что все задачи завершились
	if completedTasks != 5 {
		err := error_hndl.NewError(codes.Internal, "not all tasks completed", nil)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "registration failed: incomplete operations")
		return nil, err
	}

	// Подготавливаем запрос к базе данных
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
