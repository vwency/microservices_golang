package main

import (
	"fmt"
	"net"

	kitlog "github.com/go-kit/kit/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	"github.com/vwency/microservices_golang/internal/user_database/repository/user_repository"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	"github.com/vwency/microservices_golang/internal/user_database/transport"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/database"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "user_database", &Cfg)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(zap.NewStdLog(logger).Writer()))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	db, err := database.NewGORM(Cfg.UserDatabase.URL)
	if err != nil {
		logger.Fatal("Failed to connect to user_database", zap.Error(err))
	}

	if err := user_repository.RunUserMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize repository
	repo := repository.NewRepository(db)

	if err != nil {
		logger.Fatal("Failed to create JWT manager", zap.Error(err))
	}

	// Create service with dereferenced repository
	userService := service.NewService(*repo, kitLogger)

	// Create endpoints
	eps := endpoints.MakeEndpoints(userService)

	// Create gRPC server with both endpoints and JWT manager
	grpcServer := grpc.NewServer()
	transport.RegisterGRPCServer(grpcServer, eps)

	lis, err := net.Listen("tcp", "0.0.0.0:"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	logger.Info("Starting server",
		zap.String("service", Cfg.App.ServiceName),
		zap.String("port", Cfg.App.Port),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("gRPC server failed", zap.Error(err))
	}
}
