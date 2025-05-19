package register

import (
	"context"
	"time"

	error_hndl "github.com/vwency/microservices_golang/internal/auth_service/service/errors"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	otelAttr "go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

// validateRegisterInput проверяет обязательные поля запроса регистрации
func validateRegisterInput(ctx context.Context, tracer trace.Tracer, req *authv1.RegisterRequest) error {
	ctx, span := tracer.Start(ctx, "InputValidation")
	defer span.End()

	start := time.Now()

	if req.Username == "" || req.Password == "" || req.Email == "" {
		err := error_hndl.NewAuthError(codes.InvalidArgument, "username, password and email are required")
		span.SetAttributes(
			otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
			otelAttr.Bool("validation_passed", false),
		)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return err
	}

	span.SetAttributes(
		otelAttr.Int64("validation_duration_ns", time.Since(start).Nanoseconds()),
		otelAttr.Bool("validation_passed", true),
	)
	span.SetStatus(otelCodes.Ok, "validation passed")
	return nil
}
