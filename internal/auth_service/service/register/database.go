package register

import (
	"context"
	"time"

	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DatabaseClient interface {
	AddUser(ctx context.Context, req *databasev1.AddUserRequest, opts ...grpc.CallOption) (*databasev1.AddUserResponse, error)
}

type DatabaseService struct {
	client DatabaseClient
	tracer trace.Tracer
}

func NewDatabaseService(client DatabaseClient, tracer trace.Tracer) *DatabaseService {
	return &DatabaseService{
		client: client,
		tracer: tracer,
	}
}

func (ds *DatabaseService) AddUser(ctx context.Context, userID, username, email string, data *tokenData) error {
	ctx, span := ds.tracer.Start(ctx, "DatabaseAddUser")
	defer span.End()

	dbReq := &databasev1.AddUserRequest{
		Username:           username,
		HashedPassword:     data.hashedPassword,
		Email:              email,
		HashedAccessToken:  data.hashedAccessToken,
		HashedRefreshToken: data.hashedRefreshToken,
		UserId:             &userID,
	}

	start := time.Now()
	_, err := ds.client.AddUser(ctx, dbReq)
	span.SetAttributes(otelAttr.Int64("duration.ns", time.Since(start).Nanoseconds()))

	if err != nil {
		authError := ds.mapDatabaseError(err)
		span.RecordError(authError)
		span.SetStatus(otelCodes.Error, authError.Error())
		return authError
	}

	span.SetStatus(otelCodes.Ok, "user added to database")
	return nil
}

func (ds *DatabaseService) mapDatabaseError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.AlreadyExists:
			return error_hndl.ErrUserAlreadyExists
		default:
			return error_hndl.NewError(st.Code(), "database operation failed", st.Message())
		}
	}
	return error_hndl.NewError(codes.Internal, "database operation failed", err)
}
