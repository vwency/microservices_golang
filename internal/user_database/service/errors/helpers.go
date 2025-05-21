package errors

import (
	"errors"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target **Error) bool {
	return errors.As(err, target)
}
