package main

import (
	"fmt"
	"log"
	"net"
	"time"

	auth_service_handler "github.com/vwency/microservices_golang/internal/auth_service/handler"
	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "auth_service", &Cfg)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init zap logger: %v", err)
	}
	defer zapLogger.Sync()

	accessTokenTTL, err := time.ParseDuration(Cfg.Jwt.AccessTokenTtl)
	if err != nil {
		log.Fatalf("invalid access_token_ttl value: %v", err)
	}

	refreshTokenTTL, err := time.ParseDuration(Cfg.Jwt.RefreshTokenTtl)
	if err != nil {
		log.Fatalf("invalid refresh_token_ttl value: %v", err)
	}

	jwtManager := jwt.NewJWTManager(
		Cfg.Jwt.Secret,
		accessTokenTTL,
		refreshTokenTTL,
	)

	dbConn, err := grpc.Dial(Cfg.DatabaseService.URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to database_service: %v", err)
	}
	defer dbConn.Close()

	dbClient := databasev1.NewDatabaseInitServiceClient(dbConn)

	authUsecase := auth_service_usecase.NewAuthUsecase(dbClient, jwtManager, zapLogger)

	refreshUsecase := auth_service_usecase.NewRefreshUsecase(dbClient, jwtManager, zapLogger)
	logoutUsecase := auth_service_usecase.NewLogoutUsecase(dbClient, jwtManager, zapLogger)

	authHandler := auth_service_handler.NewServer(authUsecase, refreshUsecase, logoutUsecase, zapLogger)

	addr := fmt.Sprintf(":%s", Cfg.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", Cfg.App.Port, err)
	}

	grpcServer := grpc.NewServer()
	authHandler.RegisterService(grpcServer)

	zapLogger.Info("gRPC server for auth_service started",
		zap.String("port", Cfg.App.Port),
		zap.String("env", string(env)),
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
