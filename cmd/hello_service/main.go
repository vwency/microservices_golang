package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/hello_service/handler"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	"github.com/vwency/microservices_golang/proto/hello_service"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "hello_service", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	port := Cfg.App.Port

	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	helloHandler := handler.NewHelloHandler()
	hello_service.RegisterHelloServiceServer(grpcServer, helloHandler)

	logger.Info("gRPC server is running on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
