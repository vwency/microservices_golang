package endpoints

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

// GRPCError - кастомная ошибка с grpc кодом и сообщением
type GRPCError struct {
	Code    codes.Code
	Message string
}

func (e *GRPCError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code.String(), e.Message)
}

// GRPCErrorFromStatus создает GRPCError из gRPC статуса
func GRPCErrorFromStatus(st *status.Status) *GRPCError {
	return &GRPCError{
		Code:    st.Code(),
		Message: st.Message(),
	}
}

// Константы ошибок с grpc кодами и сообщениями
var (
	ErrInvalidArgument    = &GRPCError{Code: codes.InvalidArgument, Message: "invalid argument"}
	ErrNotFound           = &GRPCError{Code: codes.NotFound, Message: "not found"}
	ErrAlreadyExists      = &GRPCError{Code: codes.AlreadyExists, Message: "already exists"}
	ErrUnauthenticated    = &GRPCError{Code: codes.Unauthenticated, Message: "unauthenticated"}
	ErrPermissionDenied   = &GRPCError{Code: codes.PermissionDenied, Message: "permission denied"}
	ErrResourceExhausted  = &GRPCError{Code: codes.ResourceExhausted, Message: "resource exhausted"}
	ErrFailedPrecondition = &GRPCError{Code: codes.FailedPrecondition, Message: "failed precondition"}
	ErrDeadlineExceeded   = &GRPCError{Code: codes.DeadlineExceeded, Message: "deadline exceeded"}
	ErrCanceled           = &GRPCError{Code: codes.Canceled, Message: "canceled"}
	ErrUnavailable        = &GRPCError{Code: codes.Unavailable, Message: "unavailable"}
	ErrDataLoss           = &GRPCError{Code: codes.DataLoss, Message: "data loss"}
	ErrInternal           = &GRPCError{Code: codes.Internal, Message: "internal error"}
	ErrAborted            = &GRPCError{Code: codes.Aborted, Message: "operation aborted"}
)

// WrapServiceError преобразует ошибки сервисного слоя в *GRPCError с grpc кодом
func WrapServiceError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	// Если это уже gRPC статус - преобразуем в GRPCError
	if st, ok := status.FromError(err); ok {
		return GRPCErrorFromStatus(st)
	}

	// Если это уже наш GRPCError - возвращаем как есть
	var grpcErr *GRPCError
	if errors.As(err, &grpcErr) {
		return grpcErr
	}

	// Обработка ServiceError
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		switch svcErr.Code {
		case "invalid_argument":
			return ErrInvalidArgument
		case "not_found":
			return ErrNotFound
		case "already_exists":
			return ErrAlreadyExists
		case "unauthenticated":
			return ErrUnauthenticated
		case "permission_denied":
			return ErrPermissionDenied
		case "resource_exhausted":
			return ErrResourceExhausted
		case "failed_precondition":
			return ErrFailedPrecondition
		case "deadline_exceeded":
			return ErrDeadlineExceeded
		case "cancelled":
			return ErrCanceled
		case "unavailable":
			return ErrUnavailable
		case "data_loss":
			return ErrDataLoss
		case "aborted":
			return ErrAborted
		default:
			return &GRPCError{
				Code:    codes.Internal,
				Message: svcErr.Message,
			}
		}
	}

	// Обработка стандартных ошибок сервиса
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return ErrInvalidArgument
	case errors.Is(err, service.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, service.ErrAlreadyExists):
		return ErrAlreadyExists
	case errors.Is(err, service.ErrUnauthenticated):
		return ErrUnauthenticated
	case errors.Is(err, service.ErrPermissionDenied):
		return ErrPermissionDenied
	case errors.Is(err, service.ErrResourceExhausted):
		return ErrResourceExhausted
	case errors.Is(err, service.ErrFailedPrecondition):
		return ErrFailedPrecondition
	case errors.Is(err, service.ErrDeadlineExceeded):
		return ErrDeadlineExceeded
	case errors.Is(err, service.ErrCancelled):
		return ErrCanceled
	case errors.Is(err, service.ErrUnavailable):
		return ErrUnavailable
	case errors.Is(err, service.ErrDataLoss):
		return ErrDataLoss
	case errors.Is(err, service.ErrAborted):
		return ErrAborted
	case errors.Is(err, service.ErrInternal):
		return ErrInternal
	default:
		return &GRPCError{
			Code:    codes.Internal,
			Message: err.Error(),
		}
	}
}
