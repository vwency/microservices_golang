package endpoints

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError struct {
	Code    codes.Code
	Message string
}

func (e *GRPCError) Error() string {
	return e.Message
}

func (e *GRPCError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func ErrorWithDetails(code codes.Code, msg string, details ...interface{}) error {
	detailedMsg := msg
	if len(details) > 0 {
		detailedMsg = fmt.Sprintf("%s: %v", msg, details)
	}
	return &GRPCError{
		Code:    code,
		Message: detailedMsg,
	}
}

func WrapServiceError(err error) error {
	if err == nil {
		return nil
	}

	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		switch svcErr.Code {
		case "invalid_argument":
			return ErrorWithDetails(codes.InvalidArgument, svcErr.Message)
		case "not_found":
			return ErrorWithDetails(codes.NotFound, svcErr.Message)
		case "already_exists":
			return ErrorWithDetails(codes.AlreadyExists, svcErr.Message)
		case "unauthenticated":
			return ErrorWithDetails(codes.Unauthenticated, svcErr.Message)
		case "permission_denied":
			return ErrorWithDetails(codes.PermissionDenied, svcErr.Message)
		case "resource_exhausted":
			return ErrorWithDetails(codes.ResourceExhausted, svcErr.Message)
		case "failed_precondition":
			return ErrorWithDetails(codes.FailedPrecondition, svcErr.Message)
		case "deadline_exceeded":
			return ErrorWithDetails(codes.DeadlineExceeded, svcErr.Message)
		case "cancelled":
			return ErrorWithDetails(codes.Canceled, svcErr.Message)
		case "unavailable":
			return ErrorWithDetails(codes.Unavailable, svcErr.Message)
		case "data_loss":
			return ErrorWithDetails(codes.DataLoss, svcErr.Message)
		default:
			return ErrorWithDetails(codes.Internal, svcErr.Message)
		}
	}

	// Если ошибка равна какой-то из ошибок сервиса — вернем gRPC ошибку с соответствующим кодом
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return ErrorWithDetails(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return ErrorWithDetails(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrAlreadyExists):
		return ErrorWithDetails(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrUnauthenticated):
		return ErrorWithDetails(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrPermissionDenied):
		return ErrorWithDetails(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrResourceExhausted):
		return ErrorWithDetails(codes.ResourceExhausted, err.Error())
	case errors.Is(err, service.ErrFailedPrecondition):
		return ErrorWithDetails(codes.FailedPrecondition, err.Error())
	case errors.Is(err, service.ErrDeadlineExceeded):
		return ErrorWithDetails(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, service.ErrCancelled):
		return ErrorWithDetails(codes.Canceled, err.Error())
	case errors.Is(err, service.ErrUnavailable):
		return ErrorWithDetails(codes.Unavailable, err.Error())
	case errors.Is(err, service.ErrDataLoss):
		return ErrorWithDetails(codes.DataLoss, err.Error())
	default:
		return ErrorWithDetails(codes.Internal, "internal server error", err)
	}
}
