package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(meter metric.Meter) grpc.UnaryServerInterceptor {
	// Создаем счетчик
	counter, err := meter.Int64Counter("grpc.server.calls", metric.WithDescription("Total number of RPCs"))
	if err != nil {
		// Обрабатываем ошибку создания счетчика
		return nil
	}

	// Возвращаем интерсептор для отслеживания RPC
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Увеличиваем счетчик
		counter.Add(ctx, 1)

		// Вызов обработчика
		return handler(ctx, req)
	}
}
