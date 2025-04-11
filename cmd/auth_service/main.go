package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vwency/microservices_golang/internal/auth_service"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "auth_service", &Cfg)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

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
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to database_service: %v", err)
	}
	defer conn.Close()

	dbClient := databasev1.NewDatabaseInitServiceClient(conn)

	authSvc := auth_service.NewAuthService(jwtManager, logger, dbClient)

	addr := fmt.Sprintf(":%s", Cfg.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", Cfg.App.Port, err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authSvc)

	logger.Info(fmt.Sprintf("gRPC server for auth_service started on port %s", Cfg.App.Port))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
