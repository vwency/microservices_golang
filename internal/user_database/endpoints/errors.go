package endpoints

import (
	"context"
	"errors"
	"fmt"

	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
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

// ErrorWithDetails создает *GRPCError с дополнительными деталями
func ErrorWithDetails(code codes.Code, msg string, details ...interface{}) *GRPCError {
	detailedMsg := msg
	if len(details) > 0 {
		detailedMsg = fmt.Sprintf("%s: %v", msg, details)
	}
	return &GRPCError{
		Code:    code,
		Message: detailedMsg,
	}
}

// WrapServiceError преобразует ошибку сервиса или стандартную ошибку в *GRPCError
func WrapServiceError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	// Обработка ошибок контекста
	switch {
	case errors.Is(err, context.Canceled):
		return ErrorWithDetails(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return ErrorWithDetails(codes.DeadlineExceeded, "deadline exceeded")
	}

	// Если уже *GRPCError, вернуть как есть
	var grpcErr *GRPCError
	if errors.As(err, &grpcErr) {
		return grpcErr
	}

	// Если это ошибка сервиса (error_hndl.Error), конвертируем с учетом кода и сообщения
	var svcErr *error_hndl.Error
	if errors.As(err, &svcErr) {
		return ErrorWithDetails(svcErr.Code, svcErr.Message)
	}

	// Обработка ошибок из error_hndl
	switch {
	case errors.Is(err, error_hndl.ErrInvalidArgument):
		return ErrInvalidArgument
	case errors.Is(err, error_hndl.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, error_hndl.ErrAlreadyExists):
		return ErrAlreadyExists
	case errors.Is(err, error_hndl.ErrUnauthenticated):
		return ErrUnauthenticated
	case errors.Is(err, error_hndl.ErrPermissionDenied):
		return ErrPermissionDenied
	case errors.Is(err, error_hndl.ErrResourceExhausted):
		return ErrResourceExhausted
	case errors.Is(err, error_hndl.ErrFailedPrecondition):
		return ErrFailedPrecondition
	case errors.Is(err, error_hndl.ErrDeadlineExceeded):
		return ErrDeadlineExceeded
	case errors.Is(err, error_hndl.ErrCancelled):
		return ErrCanceled
	case errors.Is(err, error_hndl.ErrUnavailable):
		return ErrUnavailable
	case errors.Is(err, error_hndl.ErrDataLoss):
		return ErrDataLoss
	case errors.Is(err, error_hndl.ErrAborted):
		return ErrAborted
	case errors.Is(err, error_hndl.ErrInternal):
		return ErrInternal
	case errors.Is(err, error_hndl.ErrUnknown):
		return ErrorWithDetails(codes.Unknown, "unknown error")
	case errors.Is(err, error_hndl.ErrNotImplemented):
		return ErrorWithDetails(codes.Unimplemented, "not implemented")
	case errors.Is(err, error_hndl.ErrOutOfRange):
		return ErrorWithDetails(codes.OutOfRange, "out of range")
	}

	// По умолчанию — внутренняя ошибка с сообщением из err
	return ErrorWithDetails(codes.Internal, "internal server error", err)
}
