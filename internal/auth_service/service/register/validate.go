package register

import (
	"context"
	"regexp"
	"strings"

	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

var (
	// Регулярные выражения для валидации
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
)

// validateRegisterInput проверяет входные данные для регистрации
func validateRegisterInput(ctx context.Context, tracer trace.Tracer, req *authv1.RegisterRequest) error {
	_, span := tracer.Start(ctx, "ValidateRegisterInput",
		trace.WithAttributes(
			otelAttr.String("username", req.GetUsername()),
			otelAttr.String("email", req.GetEmail()),
			otelAttr.Int("password.length", len(req.Password)),
		))
	defer span.End()

	// Проверка username
	if err := validateUsername(req.Username); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "invalid username")
		return err
	}

	// Проверка email
	if err := validateEmail(req.Email); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "invalid email")
		return err
	}

	// Проверка password
	if err := validatePassword(req.Password); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "invalid password")
		return err
	}

	span.SetStatus(otelCodes.Ok, "validation successful")
	return nil
}

// validateUsername проверяет валидность имени пользователя
func validateUsername(username string) error {
	if len(strings.TrimSpace(username)) == 0 {
		return error_hndl.NewError(codes.InvalidArgument, "username is required", nil)
	}

	if len(username) < 3 || len(username) > 200 {
		return error_hndl.NewError(codes.InvalidArgument, "username must be between 3 and 20 characters", nil)
	}

	// if !usernameRegex.MatchString(username) {
	// 	return error_hndl.NewError(codes.InvalidArgument, "username can only contain letters, numbers, hyphens and underscores", nil)
	// }

	return nil
}

// validateEmail проверяет валидность email адреса
func validateEmail(email string) error {
	if len(strings.TrimSpace(email)) == 0 {
		return error_hndl.NewError(codes.InvalidArgument, "email is required", nil)
	}

	// if len(email) > 255 {
	// 	return error_hndl.NewError(codes.InvalidArgument, "email is too long", nil)
	// }

	// if !emailRegex.MatchString(email) {
	// 	return error_hndl.NewError(codes.InvalidArgument, "invalid email format", nil)
	// }

	return nil
}

// validatePassword проверяет валидность пароля
func validatePassword(password string) error {
	if len(password) == 0 {
		return error_hndl.NewError(codes.InvalidArgument, "password is required", nil)
	}

	// if len(password) < 8 {
	// 	return error_hndl.NewError(codes.InvalidArgument, "password must be at least 8 characters long", nil)
	// }

	// if len(password) > 128 {
	// 	return error_hndl.NewError(codes.InvalidArgument, "password is too long", nil)
	// }

	// Проверяем наличие хотя бы одной цифры и одной буквы
	// hasDigit := false
	// hasLetter := false

	// for _, char := range password {
	// 	if char >= '0' && char <= '9' {
	// 		hasDigit = true
	// 	}
	// 	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
	// 		hasLetter = true
	// 	}
	// 	if hasDigit && hasLetter {
	// 		break
	// 	}
	// }

	// if !hasDigit || !hasLetter {
	// 	return error_hndl.NewError(codes.InvalidArgument, "password must contain at least one letter and one digit", nil)
	// }

	return nil
}
