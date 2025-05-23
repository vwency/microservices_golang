package register

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"go.opentelemetry.io/otel"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

func Register(
	dbClient databasev1.DatabaseInitServiceClient,
	logger Logger,
	jwtManager JWTManager,
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

	data := tokenDataPool.Get().(*tokenData)
	defer func() {
		data.reset()
		tokenDataPool.Put(data)
	}()

	wg := wgPool.Get().(*sync.WaitGroup)
	defer wgPool.Put(wg)

	ctx, parallelSpan := tracer.Start(ctx, "ParallelOperations")
	defer parallelSpan.End()

	taskExecutor := NewTaskExecutor(
		ctx, tracer, data, wg,
		jwtManager, tokenPepper, userID,
		req.Username, req.Password,
	)

	resultChan, _, _ := taskExecutor.ExecuteAllTasks()

	completedTasks := 0
	for result := range resultChan {
		if result.err != nil {
			span.RecordError(result.err)
			span.SetStatus(otelCodes.Error, "registration failed during parallel operations")
			return nil, result.err
		}
		completedTasks++
	}

	if completedTasks != 5 {
		err := error_hndl.NewError(codes.Internal, "not all tasks completed", nil)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "registration failed: incomplete operations")
		return nil, err
	}

	parallelSpan.SetStatus(otelCodes.Ok, "all parallel operations completed successfully")

	dbService := NewDatabaseService(dbClient, tracer)
	if err := dbService.AddUser(ctx, userID, req.Username, req.Email, data); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "registration failed: database operation error")
		return nil, err
	}

	logger.Log("level", "info", "msg", "User registered successfully",
		"user_id", userID, "username", req.Username, "email", req.Email)

	span.SetStatus(otelCodes.Ok, "registration successful")

	return &authv1.RegisterResponse{
		AccessToken:  data.accessToken,
		RefreshToken: data.refreshToken,
		ExpiresAt:    data.accessExpiresAt.Unix(),
	}, nil
}
