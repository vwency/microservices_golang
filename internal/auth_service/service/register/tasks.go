package register

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	"github.com/vwency/microservices_golang/utils/authutils"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

type TaskExecutor struct {
	ctx         context.Context
	tracer      trace.Tracer
	data        *tokenData
	wg          *sync.WaitGroup
	resultChan  chan taskResult
	hasError    *int32
	semaphore   chan struct{}
	jwtManager  JWTManager
	tokenPepper string
	userID      string
	username    string
	password    string
}

type JWTManager interface {
	GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error)
	GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error)
}

func NewTaskExecutor(ctx context.Context, tracer trace.Tracer, data *tokenData, wg *sync.WaitGroup,
	jwtManager JWTManager, tokenPepper, userID, username, password string) *TaskExecutor {

	maxGoroutines := runtime.NumCPU()
	if maxGoroutines > 3 {
		maxGoroutines = 3
	}

	hasError := int32(0)

	return &TaskExecutor{
		ctx:         ctx,
		tracer:      tracer,
		data:        data,
		wg:          wg,
		resultChan:  make(chan taskResult, 5),
		hasError:    &hasError,
		semaphore:   make(chan struct{}, maxGoroutines),
		jwtManager:  jwtManager,
		tokenPepper: tokenPepper,
		userID:      userID,
		username:    username,
		password:    password,
	}
}

func (te *TaskExecutor) checkError() bool {
	select {
	case <-te.ctx.Done():
		return true
	default:
		return atomic.LoadInt32(te.hasError) != 0
	}
}

func (te *TaskExecutor) ExecuteAllTasks() (chan taskResult, chan struct{}, chan struct{}) {
	accessTokenReady := make(chan struct{})
	refreshTokenReady := make(chan struct{})

	te.wg.Add(5)

	go te.hashPasswordTask()
	go te.generateAccessTokenTask(accessTokenReady)
	go te.generateRefreshTokenTask(refreshTokenReady)
	go te.hashAccessTokenTask(accessTokenReady)
	go te.hashRefreshTokenTask(refreshTokenReady)

	go func() {
		te.wg.Wait()
		close(te.resultChan)
	}()

	return te.resultChan, accessTokenReady, refreshTokenReady
}

func (te *TaskExecutor) hashPasswordTask() {
	defer func() {
		te.wg.Done()
		<-te.semaphore
	}()
	te.semaphore <- struct{}{}

	if te.checkError() {
		te.resultChan <- taskResult{taskPassword, te.ctx.Err()}
		return
	}

	_, span := te.tracer.Start(te.ctx, "HashPassword")
	defer span.End()

	start := time.Now()
	hashedPassword, err := authutils.GenHash(te.username, te.password, nil)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		atomic.StoreInt32(te.hasError, 1)
		authErr := error_hndl.NewError(codes.Internal, "failed to hash password", err)
		span.RecordError(authErr)
		span.SetStatus(otelCodes.Error, authErr.Error())
		te.resultChan <- taskResult{taskPassword, authErr}
		return
	}

	te.data.hashedPassword = hashedPassword
	span.SetStatus(otelCodes.Ok, "password hashed successfully")
	te.resultChan <- taskResult{taskPassword, nil}
}

func (te *TaskExecutor) generateAccessTokenTask(ready chan struct{}) {
	defer func() {
		te.wg.Done()
		<-te.semaphore
	}()
	te.semaphore <- struct{}{}

	if te.checkError() {
		te.resultChan <- taskResult{taskAccessToken, te.ctx.Err()}
		return
	}

	_, span := te.tracer.Start(te.ctx, "GenerateAccessToken")
	defer span.End()

	start := time.Now()
	payload := map[string]interface{}{
		"UserID": te.userID,
		"Roles":  []interface{}{"user"},
	}

	accessToken, accessExpiresAt, err := te.jwtManager.GenerateAccessToken(payload)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		atomic.StoreInt32(te.hasError, 1)
		authErr := error_hndl.NewError(codes.Internal, "failed to generate access token", err)
		span.RecordError(authErr)
		span.SetStatus(otelCodes.Error, authErr.Error())
		te.resultChan <- taskResult{taskAccessToken, authErr}
		return
	}

	te.data.accessToken = accessToken
	te.data.accessExpiresAt = accessExpiresAt
	close(ready)
	span.SetStatus(otelCodes.Ok, "access token generated")
	te.resultChan <- taskResult{taskAccessToken, nil}
}

func (te *TaskExecutor) generateRefreshTokenTask(ready chan struct{}) {
	defer func() {
		te.wg.Done()
		<-te.semaphore
	}()
	te.semaphore <- struct{}{}

	if te.checkError() {
		te.resultChan <- taskResult{taskRefreshToken, te.ctx.Err()}
		return
	}

	_, span := te.tracer.Start(te.ctx, "GenerateRefreshToken")
	defer span.End()

	start := time.Now()
	payload := map[string]interface{}{
		"UserID": te.userID,
		"Roles":  []interface{}{"user"},
	}

	refreshToken, _, err := te.jwtManager.GenerateRefreshToken(payload)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		atomic.StoreInt32(te.hasError, 1)
		authErr := error_hndl.NewError(codes.Internal, "failed to generate refresh token", err)
		span.RecordError(authErr)
		span.SetStatus(otelCodes.Error, authErr.Error())
		te.resultChan <- taskResult{taskRefreshToken, authErr}
		return
	}

	te.data.refreshToken = refreshToken
	close(ready)
	span.SetStatus(otelCodes.Ok, "refresh token generated")
	te.resultChan <- taskResult{taskRefreshToken, nil}
}

func (te *TaskExecutor) hashAccessTokenTask(accessTokenReady chan struct{}) {
	defer func() {
		te.wg.Done()
		<-te.semaphore
	}()
	te.semaphore <- struct{}{}

	select {
	case <-te.ctx.Done():
		te.resultChan <- taskResult{taskHashAccess, te.ctx.Err()}
		return
	case <-accessTokenReady:
	}

	if te.checkError() {
		te.resultChan <- taskResult{taskHashAccess, te.ctx.Err()}
		return
	}

	_, span := te.tracer.Start(te.ctx, "HashAccessToken")
	defer span.End()

	start := time.Now()
	hashedAccessToken, err := authutils.GenHash(te.tokenPepper, te.data.accessToken, ultraFastHashParams)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		atomic.StoreInt32(te.hasError, 1)
		authErr := error_hndl.NewError(codes.Internal, "failed to hash access token", err)
		span.RecordError(authErr)
		span.SetStatus(otelCodes.Error, authErr.Error())
		te.resultChan <- taskResult{taskHashAccess, authErr}
		return
	}

	te.data.hashedAccessToken = hashedAccessToken
	span.SetStatus(otelCodes.Ok, "access token hashed")
	te.resultChan <- taskResult{taskHashAccess, nil}
}

func (te *TaskExecutor) hashRefreshTokenTask(refreshTokenReady chan struct{}) {
	defer func() {
		te.wg.Done()
		<-te.semaphore
	}()
	te.semaphore <- struct{}{}

	select {
	case <-te.ctx.Done():
		te.resultChan <- taskResult{taskHashRefresh, te.ctx.Err()}
		return
	case <-refreshTokenReady:
	}

	if te.checkError() {
		te.resultChan <- taskResult{taskHashRefresh, te.ctx.Err()}
		return
	}

	_, span := te.tracer.Start(te.ctx, "HashRefreshToken")
	defer span.End()

	start := time.Now()
	hashedRefreshToken, err := authutils.GenHash(te.tokenPepper, te.data.refreshToken, ultraFastHashParams)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		atomic.StoreInt32(te.hasError, 1)
		authErr := error_hndl.NewError(codes.Internal, "failed to hash refresh token", err)
		span.RecordError(authErr)
		span.SetStatus(otelCodes.Error, authErr.Error())
		te.resultChan <- taskResult{taskHashRefresh, authErr}
		return
	}

	te.data.hashedRefreshToken = hashedRefreshToken
	span.SetStatus(otelCodes.Ok, "refresh token hashed")
	te.resultChan <- taskResult{taskHashRefresh, nil}
}
