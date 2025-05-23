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

var (
	ErrInvalidArgument    = &GRPCError{Code: codes.InvalidArgument, Message: "invalid argument"}
	ErrNotFound           = &GRPCError{Code: codes.NotFound, Message: "not found"}
	ErrAlreadyExists      = &GRPCError{Code: codes.AlreadyExists, Message: "already exists"}
	ErrUnauthenticated    = &GRPCError{Code: codes.Unauthenticated, Message: "unauthenticated"}
	ErrPermissionDenied   = &GRPCError{Code: codes.PermissionDenied, Message: "permission denied"}
	ErrFailedPrecondition = &GRPCError{Code: codes.FailedPrecondition, Message: "failed precondition"}
	ErrResourceExhausted  = &GRPCError{Code: codes.ResourceExhausted, Message: "resource exhausted"}
	ErrInternal           = &GRPCError{Code: codes.Internal, Message: "internal error"}
	ErrUnavailable        = &GRPCError{Code: codes.Unavailable, Message: "service unavailable"}
	ErrDataLoss           = &GRPCError{Code: codes.DataLoss, Message: "data loss"}
	ErrDeadlineExceeded   = &GRPCError{Code: codes.DeadlineExceeded, Message: "deadline exceeded"}
	ErrCanceled           = &GRPCError{Code: codes.Canceled, Message: "operation cancelled"}
	ErrUnknown            = &GRPCError{Code: codes.Unknown, Message: "unknown error"}
	ErrNotImplemented     = &GRPCError{Code: codes.Unimplemented, Message: "not implemented"}
	ErrAborted            = &GRPCError{Code: codes.Aborted, Message: "operation aborted"}
	ErrOutOfRange         = &GRPCError{Code: codes.OutOfRange, Message: "out of range"}
)

func WrapServiceError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return ErrorWithDetails(codes.Canceled, "request canceled")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrorWithDetails(codes.DeadlineExceeded, "deadline exceeded")
	}

	var grpcErr *GRPCError
	if errors.As(err, &grpcErr) {
		return grpcErr
	}

	var authErr *error_hndl.Error
	if errors.As(err, &authErr) {
		return ErrorWithDetails(authErr.Code, authErr.Message)
	}

	return ErrorWithDetails(codes.Internal, "internal server error", err)
}
