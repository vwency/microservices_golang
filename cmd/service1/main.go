package main

import (
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/service1/handler"
	"github.com/vwency/microservices_golang/pkg/config" // импортируем конфиг
	"github.com/vwency/microservices_golang/pkg/logger" // импортируем логгер
	"github.com/vwency/microservices_golang/proto/service1"
)

var Cfg config.ServiceConfig // Используем тип ServiceConfig из config пакета

func main() {
	// Инициализация конфигурации
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // по умолчанию dev
	}

	// Загружаем конфигурацию для service1 из файла
	config.Init(env, "service1", &Cfg)

	// Инициализация логгера с переданным LogLevel
	logger.Init(Cfg.App.LogLevel) // передаем LogLevel из конфигурации

	// Получаем порт из конфигурации
	port := Cfg.App.Port

	// Логируем старт сервера
	logger.Info("Starting gRPC server on port " + port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	helloHandler := handler.NewHelloHandler()
	service1.RegisterHelloServiceServer(grpcServer, helloHandler)

	// Логируем, что сервер запущен
	logger.Info("gRPC server is running on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
