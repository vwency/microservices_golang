package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error представляет ошибку аутентификации с gRPC статусом
type Error struct {
	Code    codes.Code
	Message string
}

// Error реализует интерфейс error
func (e *Error) Error() string {
	return e.Message
}

// GRPCStatus реализует интерфейс GRPCStatusCoder для правильной обработки gRPC ошибок
func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// Создает новую ошибку аутентификации с деталями
func NewError(code codes.Code, msg string, details ...interface{}) *Error {
	if len(details) > 0 {
		msg = fmt.Sprintf("%s: %v", msg, details)
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// Предопределенные ошибки
var (
	ErrInvalidCredentials    = NewError(codes.Unauthenticated, "invalid credentials")
	ErrUserAlreadyExists     = NewError(codes.AlreadyExists, "user already exists")
	ErrUserNotFound          = NewError(codes.NotFound, "user not found")
	ErrInvalidToken          = NewError(codes.Unauthenticated, "invalid token")
	ErrTokenExpired          = NewError(codes.Unauthenticated, "token expired")
	ErrTokenGeneration       = NewError(codes.Internal, "failed to generate tokens")
	ErrValidationFailed      = NewError(codes.InvalidArgument, "validation failed")
	ErrPermissionDenied      = NewError(codes.PermissionDenied, "permission denied")
	ErrAccountLocked         = NewError(codes.PermissionDenied, "account is locked")
	ErrTooManyRequests       = NewError(codes.ResourceExhausted, "too many requests")
	ErrInternalServer        = NewError(codes.Internal, "internal server error")
	ErrServiceUnavailable    = NewError(codes.Unavailable, "service unavailable")
	ErrDeadlineExceeded      = NewError(codes.DeadlineExceeded, "deadline exceeded")
	ErrUnauthenticated       = NewError(codes.Unauthenticated, "unauthenticated")
	ErrRefreshTokenInvalid   = NewError(codes.Unauthenticated, "refresh token invalid")
	ErrEmailNotVerified      = NewError(codes.FailedPrecondition, "email not verified")
	ErrPasswordTooWeak       = NewError(codes.InvalidArgument, "password too weak")
	ErrUsernameTaken         = NewError(codes.AlreadyExists, "username taken")
	ErrEmailTaken            = NewError(codes.AlreadyExists, "email taken")
	ErrSessionNotFound       = NewError(codes.NotFound, "session not found")
	ErrLogoutFailed          = NewError(codes.NotFound, "logout failed")
	ErrPasswordMismatch      = NewError(codes.InvalidArgument, "password mismatch")
	ErrAccountDisabled       = NewError(codes.PermissionDenied, "account disabled")
	ErrTokenRevoked          = NewError(codes.Unauthenticated, "token revoked")
	ErrUnauthorizedOperation = NewError(codes.PermissionDenied, "unauthorized operation")
)

// Хелперы для проверки типов ошибок
func IsUserAlreadyExists(err error) bool {
	return isError(err, codes.AlreadyExists)
}

func IsInvalidCredentials(err error) bool {
	return isError(err, codes.Unauthenticated)
}

func IsUserNotFound(err error) bool {
	return isError(err, codes.NotFound)
}

func IsInvalidToken(err error) bool {
	return isError(err, codes.Unauthenticated)
}

func IsTokenExpired(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "token expired"
}

func IsTokenGeneration(err error) bool {
	return isError(err, codes.Internal)
}

func IsValidationFailed(err error) bool {
	return isError(err, codes.InvalidArgument)
}

func IsPermissionDenied(err error) bool {
	return isError(err, codes.PermissionDenied)
}

func IsAccountLocked(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "account is locked"
}

func IsTooManyRequests(err error) bool {
	return isError(err, codes.ResourceExhausted)
}

func IsInternalServer(err error) bool {
	return isError(err, codes.Internal)
}

func IsServiceUnavailable(err error) bool {
	return isError(err, codes.Unavailable)
}

func IsDeadlineExceeded(err error) bool {
	return isError(err, codes.DeadlineExceeded)
}

func IsUnauthenticated(err error) bool {
	return isError(err, codes.Unauthenticated)
}

func IsRefreshTokenInvalid(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "refresh token invalid"
}

func IsEmailNotVerified(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.FailedPrecondition && ae.Message == "email not verified"
}

func IsPasswordTooWeak(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.InvalidArgument && ae.Message == "password too weak"
}

func IsUsernameTaken(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.AlreadyExists && ae.Message == "username taken"
}

func IsEmailTaken(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.AlreadyExists && ae.Message == "email taken"
}

func IsSessionNotFound(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.NotFound && ae.Message == "session not found"
}

func IsLogoutFailed(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.Internal && ae.Message == "logout failed"
}

func IsPasswordMismatch(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.InvalidArgument && ae.Message == "password mismatch"
}

func IsAccountDisabled(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "account disabled"
}

func IsTokenRevoked(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "token revoked"
}

func IsUnauthorizedOperation(err error) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "unauthorized operation"
}

func isError(err error, code codes.Code) bool {
	ae, ok := err.(*Error)
	return ok && ae.Code == code
}
