package errors

import (
	"google.golang.org/grpc/codes"
)

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
