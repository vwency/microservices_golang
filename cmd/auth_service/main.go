package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vwency/microservices_golang/internal/auth_service"
	"github.com/vwency/microservices_golang/internal/database/repository"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "auth_service", &Cfg)

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
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

	// ✅ Подключение к БД, используя строку подключения из Cfg.Database.URL
	db, err := gorm.Open(postgres.Open(Cfg.Database.URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// ✅ Создание репозитория и usecase
	repo := repository.NewRepository(db)
	userUC := user_usecase.New(repo.UserRepo, logger)

	addr := fmt.Sprintf(":%s", Cfg.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", Cfg.App.Port, err)
	}

	grpcServer := grpc.NewServer()

	authSvc := auth_service.NewAuthService(jwtManager, userUC)
	authv1.RegisterAuthServiceServer(grpcServer, authSvc)

	logger.Info(fmt.Sprintf("gRPC server for auth_service started on port %s", Cfg.App.Port))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
