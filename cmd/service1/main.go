package main

import (
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/service1/handler"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	"github.com/vwency/microservices_golang/proto/service1"
)

var Cfg config.ServiceConfig

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	config.Init(env, "service1", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	port := Cfg.App.Port

	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	helloHandler := handler.NewHelloHandler()
	service1.RegisterHelloServiceServer(grpcServer, helloHandler)

	logger.Info("gRPC server is running on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
