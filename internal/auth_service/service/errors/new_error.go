package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code    codes.Code
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewError(code codes.Code, msg string, details ...interface{}) *Error {
	if len(details) > 0 {
		msg = fmt.Sprintf("%s: %v", msg, details)
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}
