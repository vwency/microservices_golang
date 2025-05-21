package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthError представляет ошибку аутентификации с gRPC статусом
type AuthError struct {
	Code    codes.Code
	Message string
}

// Error реализует интерфейс error
func (e *AuthError) Error() string {
	return e.Message
}

// GRPCStatus реализует интерфейс GRPCStatusCoder для правильной обработки gRPC ошибок
func (e *AuthError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// Создает новую ошибку аутентификации с деталями
func NewAuthError(code codes.Code, msg string, details ...interface{}) *AuthError {
	if len(details) > 0 {
		msg = fmt.Sprintf("%s: %v", msg, details)
	}
	return &AuthError{
		Code:    code,
		Message: msg,
	}
}

// Предопределенные ошибки
var (
	ErrInvalidCredentials    = NewAuthError(codes.Unauthenticated, "invalid credentials")
	ErrUserAlreadyExists     = NewAuthError(codes.AlreadyExists, "user already exists")
	ErrUserNotFound          = NewAuthError(codes.NotFound, "user not found")
	ErrInvalidToken          = NewAuthError(codes.Unauthenticated, "invalid token")
	ErrTokenExpired          = NewAuthError(codes.Unauthenticated, "token expired")
	ErrTokenGeneration       = NewAuthError(codes.Internal, "failed to generate tokens")
	ErrValidationFailed      = NewAuthError(codes.InvalidArgument, "validation failed")
	ErrPermissionDenied      = NewAuthError(codes.PermissionDenied, "permission denied")
	ErrAccountLocked         = NewAuthError(codes.PermissionDenied, "account is locked")
	ErrTooManyRequests       = NewAuthError(codes.ResourceExhausted, "too many requests")
	ErrInternalServer        = NewAuthError(codes.Internal, "internal server error")
	ErrServiceUnavailable    = NewAuthError(codes.Unavailable, "service unavailable")
	ErrDeadlineExceeded      = NewAuthError(codes.DeadlineExceeded, "deadline exceeded")
	ErrUnauthenticated       = NewAuthError(codes.Unauthenticated, "unauthenticated")
	ErrRefreshTokenInvalid   = NewAuthError(codes.Unauthenticated, "refresh token invalid")
	ErrEmailNotVerified      = NewAuthError(codes.FailedPrecondition, "email not verified")
	ErrPasswordTooWeak       = NewAuthError(codes.InvalidArgument, "password too weak")
	ErrUsernameTaken         = NewAuthError(codes.AlreadyExists, "username taken")
	ErrEmailTaken            = NewAuthError(codes.AlreadyExists, "email taken")
	ErrSessionNotFound       = NewAuthError(codes.NotFound, "session not found")
	ErrLogoutFailed          = NewAuthError(codes.NotFound, "logout failed")
	ErrPasswordMismatch      = NewAuthError(codes.InvalidArgument, "password mismatch")
	ErrAccountDisabled       = NewAuthError(codes.PermissionDenied, "account disabled")
	ErrTokenRevoked          = NewAuthError(codes.Unauthenticated, "token revoked")
	ErrUnauthorizedOperation = NewAuthError(codes.PermissionDenied, "unauthorized operation")
)

// Хелперы для проверки типов ошибок
func IsUserAlreadyExists(err error) bool {
	return isAuthError(err, codes.AlreadyExists)
}

func IsInvalidCredentials(err error) bool {
	return isAuthError(err, codes.Unauthenticated)
}

func IsUserNotFound(err error) bool {
	return isAuthError(err, codes.NotFound)
}

func IsInvalidToken(err error) bool {
	return isAuthError(err, codes.Unauthenticated)
}

func IsTokenExpired(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "token expired"
}

func IsTokenGeneration(err error) bool {
	return isAuthError(err, codes.Internal)
}

func IsValidationFailed(err error) bool {
	return isAuthError(err, codes.InvalidArgument)
}

func IsPermissionDenied(err error) bool {
	return isAuthError(err, codes.PermissionDenied)
}

func IsAccountLocked(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "account is locked"
}

func IsTooManyRequests(err error) bool {
	return isAuthError(err, codes.ResourceExhausted)
}

func IsInternalServer(err error) bool {
	return isAuthError(err, codes.Internal)
}

func IsServiceUnavailable(err error) bool {
	return isAuthError(err, codes.Unavailable)
}

func IsDeadlineExceeded(err error) bool {
	return isAuthError(err, codes.DeadlineExceeded)
}

func IsUnauthenticated(err error) bool {
	return isAuthError(err, codes.Unauthenticated)
}

func IsRefreshTokenInvalid(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "refresh token invalid"
}

func IsEmailNotVerified(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.FailedPrecondition && ae.Message == "email not verified"
}

func IsPasswordTooWeak(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.InvalidArgument && ae.Message == "password too weak"
}

func IsUsernameTaken(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.AlreadyExists && ae.Message == "username taken"
}

func IsEmailTaken(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.AlreadyExists && ae.Message == "email taken"
}

func IsSessionNotFound(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.NotFound && ae.Message == "session not found"
}

func IsLogoutFailed(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.Internal && ae.Message == "logout failed"
}

func IsPasswordMismatch(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.InvalidArgument && ae.Message == "password mismatch"
}

func IsAccountDisabled(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "account disabled"
}

func IsTokenRevoked(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.Unauthenticated && ae.Message == "token revoked"
}

func IsUnauthorizedOperation(err error) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == codes.PermissionDenied && ae.Message == "unauthorized operation"
}

func isAuthError(err error, code codes.Code) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == code
}
