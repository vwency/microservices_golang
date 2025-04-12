package metrics

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func NewMeterProvider(cfg Config) (*metric.MeterProvider, error) {
	if !cfg.EnableMetrics {
		return nil, nil
	}

	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OtlpEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("Failed to create OTLP metrics exporter: %v", err)
		return nil, err
	}

	if cfg.ExportInterval == 0 {
		cfg.ExportInterval = 10 * time.Second
	}
	if cfg.ExportTimeout == 0 {
		cfg.ExportTimeout = 5 * time.Second
	}

	reader := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(cfg.ExportInterval),
		metric.WithTimeout(cfg.ExportTimeout),
	)

	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}
