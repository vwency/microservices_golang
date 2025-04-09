package main

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	handler_user_service "github.com/vwency/microservices_golang/internal/database/delivery/grpc/handler/user_service"
	"github.com/vwency/microservices_golang/internal/database/repository"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/database"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "database_service", &Cfg)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	db, err := database.NewGORM(Cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	repo := repository.NewRepository(db)
	userUC := user_usecase.New(repo.UserRepo, logger)

	lis, err := net.Listen("tcp", "0.0.0.0:"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	userHandler := handler_user_service.NewServer(userUC, logger)
	userHandler.Register(grpcServer)

	logger.Info("Starting server",
		zap.String("service", Cfg.App.ServiceName),
		zap.String("port", Cfg.App.Port),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("gRPC server failed", zap.Error(err))
	}
}
