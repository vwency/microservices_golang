package main

import (
	"net"

	handler_user_service "github.com/vwency/microservices_golang/internal/database/delivery/gprc/handler/user_service"
	"github.com/vwency/microservices_golang/internal/database/repository"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase" // исправляем импорт
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	pb "github.com/vwency/microservices_golang/proto/database"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "database_service", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	dsn := Cfg.Database.URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	repo := repository.NewRepository(db)

	uc := user_usecase.NewInitUseCase(repo.UserRepo)

	server := handler_user_service.NewServer(uc)

	lis, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterDatabaseInitServiceServer(grpcServer, server)

	logger.Info("Starting database init service on port " + Cfg.App.Port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
