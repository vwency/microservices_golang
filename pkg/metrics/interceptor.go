package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(meter metric.Meter) grpc.UnaryServerInterceptor {
	counter, err := meter.Int64Counter("grpc.server.calls", metric.WithDescription("Total number of RPCs"))
	if err != nil {
		return nil
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		counter.Add(ctx, 1)

		return handler(ctx, req)
	}
}
