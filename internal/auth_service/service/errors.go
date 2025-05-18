package service

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
	ErrInvalidCredentials = NewAuthError(codes.Unauthenticated, "invalid credentials")
	ErrUserAlreadyExists  = NewAuthError(codes.AlreadyExists, "user already exists")
	ErrUserNotFound       = NewAuthError(codes.NotFound, "user not found")
	ErrInvalidToken       = NewAuthError(codes.Unauthenticated, "invalid token")
	ErrTokenExpired       = NewAuthError(codes.Unauthenticated, "token expired")
	ErrTokenGeneration    = NewAuthError(codes.Internal, "failed to generate tokens")
	ErrValidationFailed   = NewAuthError(codes.InvalidArgument, "validation failed")
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

func isAuthError(err error, code codes.Code) bool {
	ae, ok := err.(*AuthError)
	return ok && ae.Code == code
}
