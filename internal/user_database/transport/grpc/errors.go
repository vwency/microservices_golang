package grpc

import (
	"context"
	"errors"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Уже gRPC статус?
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Контекстные ошибки
	switch {
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	}

	// Разворачиваем ошибку
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return ConvertToGRPCError(unwrapped)
	}

	// Ошибки endpoints
	switch {
	case errors.Is(err, endpoints.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, endpoints.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, endpoints.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, endpoints.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, endpoints.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, endpoints.ErrResourceExhausted):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, endpoints.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, endpoints.ErrDeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, endpoints.ErrCanceled):
		return status.Error(codes.Canceled, err.Error())
	case errors.Is(err, endpoints.ErrUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, endpoints.ErrDataLoss):
		return status.Error(codes.DataLoss, err.Error())
	}

	// Ошибки сервиса напрямую
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		switch svcErr.Code {
		case "invalid_argument":
			return status.Error(codes.InvalidArgument, svcErr.Message)
		case "not_found":
			return status.Error(codes.NotFound, svcErr.Message)
		case "already_exists":
			return status.Error(codes.AlreadyExists, svcErr.Message)
		case "unauthenticated":
			return status.Error(codes.Unauthenticated, svcErr.Message)
		case "permission_denied":
			return status.Error(codes.PermissionDenied, svcErr.Message)
		case "resource_exhausted":
			return status.Error(codes.ResourceExhausted, svcErr.Message)
		case "failed_precondition":
			return status.Error(codes.FailedPrecondition, svcErr.Message)
		case "deadline_exceeded":
			return status.Error(codes.DeadlineExceeded, svcErr.Message)
		case "cancelled":
			return status.Error(codes.Canceled, svcErr.Message)
		case "unavailable":
			return status.Error(codes.Unavailable, svcErr.Message)
		case "data_loss":
			return status.Error(codes.DataLoss, svcErr.Message)
		default:
			return status.Error(codes.Internal, svcErr.Message)
		}
	}

	// По умолчанию
	return status.Error(codes.Internal, "internal server error: "+err.Error())
}
