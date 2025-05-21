package endpoints

import (
	"context"
	"errors"
	"fmt"

	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
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
	ErrInvalidArgument = &GRPCError{Code: codes.InvalidArgument, Message: "invalid argument"}
	ErrNotFound        = &GRPCError{Code: codes.NotFound, Message: "not found"}
)

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

func WrapServiceError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return ErrorWithDetails(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
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
