package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vwency/microservices_golang/internal/auth_service"
	"github.com/vwency/microservices_golang/pkg/auth"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"

	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	// Детектируем окружение и загружаем конфиг
	env := config.DetectEnv()
	config.Init(env, "auth_service", &Cfg)

	// Инициализация логгера
	logger.Init(Cfg.App.LogLevel)

	// Преобразуем TTL из строк в time.Duration
	accessTokenTTL, err := time.ParseDuration(Cfg.Jwt.AccessTokenTtl)
	if err != nil {
		log.Fatalf("invalid access_token_ttl value: %v", err)
	}

	refreshTokenTTL, err := time.ParseDuration(Cfg.Jwt.RefreshTokenTtl)
	if err != nil {
		log.Fatalf("invalid refresh_token_ttl value: %v", err)
	}

	// Создание JWT менеджера
	jwtManager := auth.NewJWTManager(
		Cfg.Jwt.Secret,
		accessTokenTTL,
		refreshTokenTTL,
	)

	// Настройка TCP-порта
	addr := fmt.Sprintf(":%s", Cfg.App.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", Cfg.App.Port, err)
	}

	grpcServer := grpc.NewServer()

	// Создание и регистрация AuthService
	authSvc := auth_service.NewAuthService(jwtManager)
	authv1.RegisterAuthServiceServer(grpcServer, authSvc)

	logger.Info(fmt.Sprintf("gRPC server for auth_service started on port %s", Cfg.App.Port))

	// Запуск gRPC-сервера
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
