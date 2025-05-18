package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeRegisterEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Создаем tracer
		tracer := otel.Tracer("auth_service.endpoint")

		// Span для всей обработки endpoint
		ctx, span := tracer.Start(ctx, "RegisterEndpoint",
			trace.WithAttributes(
				attribute.String("component", "endpoint"),
				attribute.String("span.kind", "server"),
			))
		defer span.End()

		// 1. Проверка типа запроса
		_, typeCheckSpan := tracer.Start(ctx, "Register.TypeCheck")
		req, ok := request.(*authv1.RegisterRequest)
		typeCheckSpan.End()

		if !ok {
			span.RecordError(status.Error(grpc_codes.InvalidArgument, "invalid request type"))
			span.SetStatus(codes.Error, "invalid request type")
			return nil, status.Error(grpc_codes.InvalidArgument, "invalid request type")
		}

		// Добавляем атрибуты запроса в span
		span.SetAttributes(
			attribute.String("request.email", req.GetEmail()),
			attribute.String("request.username", req.GetUsername()),
		)

		// 2. Вызов сервисного слоя
		_, serviceCallSpan := tracer.Start(ctx, "Register.ServiceCall")
		res, err := s.Register(ctx, req)
		serviceCallSpan.End()

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		// Добавляем атрибуты ответа
		if res != nil {
			span.SetAttributes(
				attribute.String("response.access_token", res.GetAccessToken()),
				attribute.Int64("response.expires_at", res.GetExpiresAt()),
			)
		}

		span.SetStatus(codes.Ok, "success")
		return res, nil
	}
}
