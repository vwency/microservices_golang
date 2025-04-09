package main

import (
	"net"

	"github.com/vwency/microservices_golang/internal/database/delivery/gprc/handler"
	"github.com/vwency/microservices_golang/internal/database/repository"
	"github.com/vwency/microservices_golang/internal/database/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	pb "github.com/vwency/microservices_golang/proto/database"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Cfg config.ServiceConfig

func main() {
	// Initialize configuration
	env := config.DetectEnv()
	config.Init(env, "database_init_service", &Cfg)

	// Initialize logger
	logger.Init(Cfg.App.LogLevel)

	// Connect to database
	dsn := Cfg.Database.URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	// Initialize repository
	repo := repository.NewUserRepository(db) // This is now passed correctly as an interface

	// Initialize UseCase
	uc := usecase.NewInitUseCase(repo)

	// Initialize handler
	handler := handler.NewHandler(uc) // Make sure we are calling NewHandler

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the service
	pb.RegisterDatabaseInitServiceServer(grpcServer, handler)

	logger.Info("Starting database init service on port " + Cfg.App.Port)

	// Start gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
