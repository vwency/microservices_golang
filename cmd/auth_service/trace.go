package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
)

func registerTracerShutdown(lc fx.Lifecycle, tp *sdktrace.TracerProvider) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// graceful shutdown TracerProvider (отправка всех спанов)
			return tp.Shutdown(ctx)
		},
	})
}

func initTrace(lc fx.Lifecycle) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("127.0.0.1:4317"),
		// Remove WithBlock to prevent hanging if collector isn't ready
		otlptracegrpc.WithReconnectionPeriod(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create OTLP trace exporter: %v", err)
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("auth_service"), // Changed to match your service
			semconv.ServiceVersion("1.0.0"),
		)),
	)

	// Add shutdown hook
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})

	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Printf("OpenTelemetry error: %v", err)
	}))

	log.Println("OpenTelemetry tracer initialized successfully")
	return tp, nil
}
