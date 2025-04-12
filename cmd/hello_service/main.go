package main

import (
	"context"
	"net"
	"time"

	handler_hello "github.com/vwency/microservices_golang/internal/hello_service/handler"
	usecase_hello "github.com/vwency/microservices_golang/internal/hello_service/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	"github.com/vwency/microservices_golang/pkg/metrics"
	"github.com/vwency/microservices_golang/pkg/tracing"
	"github.com/vwency/microservices_golang/proto/hello_service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "hello_service", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	tp, err := tracing.NewTracerProvider(tracing.Config{
		ServiceName:   Cfg.App.ServiceName,
		EnableTracing: Cfg.Tracing.Enabled,
		OtlpEndpoint:  Cfg.Tracing.OtlpEndpoint,
	})
	if err != nil {
		logger.Fatal("failed to initialize tracing: %v", err)
	}
	if tp != nil {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logger.Error("failed to shutdown tracer provider: %v", err)
			}
		}()
	}

	// Инициализация метрик
	var meter metric.Meter
	if Cfg.Metrics.Enabled {
		mp, err := metrics.NewMeterProvider(metrics.Config{
			ServiceName:    Cfg.App.ServiceName,
			EnableMetrics:  Cfg.Metrics.Enabled,
			OtlpEndpoint:   Cfg.Metrics.OtlpEndpoint,
			ExportInterval: parseDurationOrDefault(Cfg.Metrics.ExportInterval, 10*time.Second),
			ExportTimeout:  parseDurationOrDefault(Cfg.Metrics.ExportTimeout, 5*time.Second),
		})
		if err != nil {
			logger.Fatal("failed to initialize metrics: %v", err)
		}
		if mp != nil {
			defer func() {
				if err := mp.Shutdown(context.Background()); err != nil {
					logger.Error("failed to shutdown meter provider: %v", err)
				}
			}()
			// Get the meter from the global provider (since we called otel.SetMeterProvider)
			meter = otel.Meter(Cfg.App.ServiceName)
		}
	}

	port := Cfg.App.Port
	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	// Создаем цепочку интерсепторов
	interceptors := []grpc.UnaryServerInterceptor{}

	if tp != nil {
		interceptors = append(interceptors, tracing.TracingInterceptor(tp.Tracer(Cfg.App.ServiceName)))
	}

	if meter != nil {
		interceptors = append(interceptors, metrics.MetricsInterceptor(meter))
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	helloUsecase := usecase_hello.NewHelloUsecase()
	helloHandler := handler_hello.NewHelloHandler(helloUsecase)

	hello_service.RegisterHelloServiceServer(grpcServer, helloHandler)

	logger.Info("gRPC server is running on port " + port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}

// parseDurationOrDefault parses a duration string or returns a default value
func parseDurationOrDefault(durationStr string, defaultValue time.Duration) time.Duration {
	if durationStr == "" {
		return defaultValue
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Debug("invalid duration %q, using default: %v", durationStr, err)
		return defaultValue
	}
	return dur
}
