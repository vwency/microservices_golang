package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument = NewGRPCError(codes.InvalidArgument, "invalid argument")
	ErrNotFound        = NewGRPCError(codes.NotFound, "not found")
	ErrInternal        = NewGRPCError(codes.Internal, "internal error")
	ErrAlreadyExists   = NewGRPCError(codes.AlreadyExists, "already exists")
)

type GRPCError struct {
	code    codes.Code
	message string
}

func NewGRPCError(code codes.Code, message string) *GRPCError {
	return &GRPCError{code: code, message: message}
}

func (e *GRPCError) Error() string {
	return e.message
}

func (e *GRPCError) GRPCStatus() *status.Status {
	return status.New(e.code, e.message)
}
