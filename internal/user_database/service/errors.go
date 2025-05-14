package service

import (
	"errors"
)

var (
	// 4xx errors
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrResourceExhausted  = errors.New("resource exhausted")

	// 5xx errors
	ErrInternal         = errors.New("internal error")
	ErrUnavailable      = errors.New("service unavailable")
	ErrDataLoss         = errors.New("data loss")
	ErrDeadlineExceeded = errors.New("deadline exceeded")

	// Other
	ErrCancelled      = errors.New("operation cancelled")
	ErrUnknown        = errors.New("unknown error")
	ErrNotImplemented = errors.New("not implemented")
	ErrAborted        = errors.New("operation aborted")
	ErrOutOfRange     = errors.New("out of range")
)

type ServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

func NewInvalidArgumentError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "invalid_argument",
		Message: msg,
		Err:     err,
	}
}

func NewNotFoundError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "not_found",
		Message: msg,
		Err:     err,
	}
}

func NewAlreadyExistsError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "already_exists",
		Message: msg,
		Err:     err,
	}
}

func NewUnauthenticatedError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "unauthenticated",
		Message: msg,
		Err:     err,
	}
}

func NewPermissionDeniedError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "permission_denied",
		Message: msg,
		Err:     err,
	}
}

func NewInternalError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "internal",
		Message: msg,
		Err:     err,
	}
}

func NewResourceExhaustedError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "resource_exhausted",
		Message: msg,
		Err:     err,
	}
}

func NewFailedPreconditionError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "failed_precondition",
		Message: msg,
		Err:     err,
	}
}

func NewUnavailableError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "unavailable",
		Message: msg,
		Err:     err,
	}
}

func NewDataLossError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "data_loss",
		Message: msg,
		Err:     err,
	}
}

func NewDeadlineExceededError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "deadline_exceeded",
		Message: msg,
		Err:     err,
	}
}

func NewCancelledError(msg string, err error) *ServiceError {
	return &ServiceError{
		Code:    "cancelled",
		Message: msg,
		Err:     err,
	}
}
