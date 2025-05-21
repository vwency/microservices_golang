package errors

import (
	"errors"
)

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrResourceExhausted  = errors.New("resource exhausted")

	ErrInternal         = errors.New("internal error")
	ErrUnavailable      = errors.New("service unavailable")
	ErrDataLoss         = errors.New("data loss")
	ErrDeadlineExceeded = errors.New("deadline exceeded")

	ErrCancelled      = errors.New("operation cancelled")
	ErrUnknown        = errors.New("unknown error")
	ErrNotImplemented = errors.New("not implemented")
	ErrAborted        = errors.New("operation aborted")
	ErrOutOfRange     = errors.New("out of range")
)
