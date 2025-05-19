package endpoints

import (
	"errors"

	"github.com/vwency/microservices_golang/internal/user_database/service"
)

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrResourceExhausted  = errors.New("resource exhausted")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrDeadlineExceeded   = errors.New("deadline exceeded")
	ErrCanceled           = errors.New("canceled")
	ErrUnavailable        = errors.New("unavailable")
	ErrDataLoss           = errors.New("data loss")
	ErrInternal           = errors.New("internal error")
)

func WrapServiceError(err error) error {
	if err == nil {
		return nil
	}

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
		default:
			return ErrInternal
		}
	}

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
	case errors.Is(err, service.ErrInternal):
		return ErrInternal
	default:
		return ErrInternal
	}
}
