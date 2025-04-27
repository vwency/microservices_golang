package main

import (
	"net"
	"time"

	handler_jwt "github.com/vwency/microservices_golang/internal/jwt_service/handler"
	usecase_jwt "github.com/vwency/microservices_golang/internal/jwt_service/usecase"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
	"github.com/vwency/microservices_golang/proto/jwt_service"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	// Инициализация конфигурации
	env := config.DetectEnv()
	config.Init(env, "jwt_service", &Cfg)

	// Инициализация логгера
	logger.Init(Cfg.App.LogLevel)

	// Инициализация gRPC сервера
	port := Cfg.App.Port
	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	// Интерсепторы для gRPC
	interceptors := []grpc.UnaryServerInterceptor{}

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	// Инициализация Usecase и Handler для jwt_service
	jwtUsecase := usecase_jwt.NewJwtUsecase()
	jwtHandler := handler_jwt.NewJwtHandler(jwtUsecase)

	// Регистрация сервиса в gRPC сервере
	jwt_service.RegisterJwtServiceServer(grpcServer, jwtHandler)

	// Запуск сервера
	logger.Info("gRPC server is running on port " + port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}

// Утилита для парсинга длительности с дефолтным значением
func parseDurationOrDefault(durationStr string, defaultValue time.Duration) time.Duration {
	if durationStr == "" {
		return defaultValue
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Debug("invalid duration %q, using default: %v", durationStr, err)
		return defaultValue
	}
	return dur
}
