package endpoints

import (
	"context"
	"log"
	"time"

	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TraceEndpoint создает middleware для трассировки вызовов endpoint
func TraceEndpoint(tracer trace.Tracer, name string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// Начинаем замер времени выполнения
			startTime := time.Now()

			// Создаем span с учетом возможного родительского span из контекста
			ctx, span := tracer.Start(
				trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
				name,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("endpoint.name", name),
					attribute.String("component", "endpoint"),
				),
			)
			defer func() {
				// Завершаем span и логируем результат
				duration := time.Since(startTime)
				span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))

				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					log.Printf("[TraceEndpoint] %s failed after %s: %v", name, duration, err)
				} else {
					span.SetStatus(codes.Ok, "success")
					log.Printf("[TraceEndpoint] %s completed in %s", name, duration)
				}

				span.End()
			}()

			if span == nil {
				log.Printf("[TraceEndpoint] Failed to create span for %s", name)
				return next(ctx, request)
			}

			// Добавляем атрибуты из запроса
			if req, ok := request.(interface{ GetTraceInfo() map[string]string }); ok && req != nil {
				for k, v := range req.GetTraceInfo() {
					if v != "" { // Добавляем только непустые значения
						span.SetAttributes(attribute.String("request."+k, v))
					}
				}
			}

			// Выполняем следующий обработчик в цепочке
			response, err = next(ctx, request)

			// Добавляем атрибуты из ответа
			if err == nil && response != nil {
				if resp, ok := response.(interface{ GetTraceInfo() map[string]string }); ok && resp != nil {
					for k, v := range resp.GetTraceInfo() {
						if v != "" {
							span.SetAttributes(attribute.String("response."+k, v))
						}
					}
				}
			}

			return response, err
		}
	}
}
