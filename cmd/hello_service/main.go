package main

import (
	"context"
	"net"

	handler_hello "github.com/vwency/microservices_golang/internal/hello_service/handler"
	usecase_hello "github.com/vwency/microservices_golang/internal/hello_service/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	"github.com/vwency/microservices_golang/pkg/tracing"
	"github.com/vwency/microservices_golang/proto/hello_service"
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
			_ = tp.Shutdown(context.Background())
		}()
	}

	port := Cfg.App.Port
	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(tracing.TracingInterceptor(tp.Tracer(Cfg.App.ServiceName))),
	)

	helloUsecase := usecase_hello.NewHelloUsecase()
	helloHandler := handler_hello.NewHelloHandler(helloUsecase)

	hello_service.RegisterHelloServiceServer(grpcServer, helloHandler)

	logger.Info("gRPC server is running on port " + port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
