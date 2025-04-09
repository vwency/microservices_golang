package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/database/handler"
	repository "github.com/vwency/microservices_golang/internal/database/repository/user"
	"github.com/vwency/microservices_golang/internal/database/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	pb "github.com/vwency/microservices_golang/proto/database"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "database_init_service", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	repo, err := repository.NewUserRepository(Cfg)
	if err != nil {
		logger.Fatal("failed to initialize repository: %v", err)
	}

	uc := usecase.NewInitUseCase(repo)

	handler := handler.NewHandler(uc)

	lis, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterDatabaseInitServiceServer(grpcServer, handler)

	logger.Info("Starting database init service on port " + Cfg.App.Port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
