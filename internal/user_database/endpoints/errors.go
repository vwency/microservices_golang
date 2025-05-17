package endpoints

import (
	"errors"
	"fmt"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCError represents a gRPC error with status code and message
type GRPCError struct {
	Code    codes.Code
	Message string
}

// Error implements the error interface
func (e *GRPCError) Error() string {
	return e.Message
}

// GRPCStatus implements the GRPCStatusCoder interface for proper gRPC error handling
func (e *GRPCError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// Predefined gRPC errors
var (
	ErrInvalidArgument    = &GRPCError{Code: codes.InvalidArgument, Message: "invalid argument"}
	ErrNotFound           = &GRPCError{Code: codes.NotFound, Message: "not found"}
	ErrInternal           = &GRPCError{Code: codes.Internal, Message: "internal error"}
	ErrAlreadyExists      = &GRPCError{Code: codes.AlreadyExists, Message: "already exists"}
	ErrUnauthenticated    = &GRPCError{Code: codes.Unauthenticated, Message: "unauthenticated"}
	ErrPermissionDenied   = &GRPCError{Code: codes.PermissionDenied, Message: "permission denied"}
	ErrResourceExhausted  = &GRPCError{Code: codes.ResourceExhausted, Message: "resource exhausted"}
	ErrFailedPrecondition = &GRPCError{Code: codes.FailedPrecondition, Message: "failed precondition"}
	ErrDeadlineExceeded   = &GRPCError{Code: codes.DeadlineExceeded, Message: "deadline exceeded"}
	ErrCancelled          = &GRPCError{Code: codes.Canceled, Message: "request cancelled"}
	ErrUnavailable        = &GRPCError{Code: codes.Unavailable, Message: "service unavailable"}
	ErrDataLoss           = &GRPCError{Code: codes.DataLoss, Message: "data loss"}
)

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

	if errors.Is(err, service.ErrInvalidArgument) {
		return ErrInvalidArgument
	}
	if errors.Is(err, service.ErrNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, service.ErrAlreadyExists) {
		return ErrAlreadyExists
	}
	if errors.Is(err, service.ErrUnauthenticated) {
		return ErrUnauthenticated
	}
	if errors.Is(err, service.ErrPermissionDenied) {
		return ErrPermissionDenied
	}
	if errors.Is(err, service.ErrResourceExhausted) {
		return ErrResourceExhausted
	}
	if errors.Is(err, service.ErrFailedPrecondition) {
		return ErrFailedPrecondition
	}
	if errors.Is(err, service.ErrDeadlineExceeded) {
		return ErrDeadlineExceeded
	}
	if errors.Is(err, service.ErrCancelled) {
		return ErrCancelled
	}
	if errors.Is(err, service.ErrUnavailable) {
		return ErrUnavailable
	}
	if errors.Is(err, service.ErrDataLoss) {
		return ErrDataLoss
	}

	return ErrorWithDetails(codes.Internal, "internal server error", err)
}
