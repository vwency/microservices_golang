package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
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

// ErrorWithDetails creates a new gRPC error with additional details
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

// ErrorWrapperMiddleware wraps service errors into appropriate gRPC status errors
func ErrorWrapperMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			resp, err := next(ctx, request)
			if err == nil {
				return resp, nil
			}

			logger.Log("error", err)

			// If error already has gRPC status, return as is
			if st, ok := status.FromError(err); ok {
				return nil, st.Err()
			}

			// Handle service errors
			switch {
			case error_hndl.IsUserAlreadyExists(err):
				return nil, ErrorWithDetails(codes.AlreadyExists, err.Error())
			case error_hndl.IsInvalidCredentials(err):
				return nil, ErrorWithDetails(codes.Unauthenticated, err.Error())
			case error_hndl.IsUserNotFound(err):
				return nil, ErrorWithDetails(codes.NotFound, err.Error())
			case error_hndl.IsInvalidToken(err) || error_hndl.IsTokenExpired(err):
				return nil, ErrorWithDetails(codes.Unauthenticated, err.Error())
			case error_hndl.IsTokenGeneration(err):
				return nil, ErrorWithDetails(codes.Internal, err.Error())
			case error_hndl.IsValidationFailed(err):
				return nil, ErrorWithDetails(codes.InvalidArgument, err.Error())
			default:
				// Wrap generic errors
				return nil, WrapServiceError(err)
			}
		}
	}
}

// WrapServiceError wraps service errors into gRPC errors
func WrapServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Handle context errors
	switch {
	case errors.Is(err, context.Canceled):
		return ErrorWithDetails(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return ErrorWithDetails(codes.DeadlineExceeded, "deadline exceeded")
	}

	// If error is already a GRPCError, return as is
	var grpcErr *GRPCError
	if errors.As(err, &grpcErr) {
		return grpcErr
	}

	// If error is an AuthError, convert to GRPCError
	var authErr *error_hndl.Error
	if errors.As(err, &authErr) {
		return ErrorWithDetails(authErr.Code, authErr.Message)
	}

	// Default case - internal server error
	return ErrorWithDetails(codes.Internal, "internal server error", err)
}
