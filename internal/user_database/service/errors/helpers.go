package errors

import (
	"google.golang.org/grpc/codes"
)

func isError(err error, code codes.Code) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == code
}
