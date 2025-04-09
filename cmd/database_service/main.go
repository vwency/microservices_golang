package main

import (
	"net"

	"github.com/vwency/microservices_golang/internal/database/handler"
	repository "github.com/vwency/microservices_golang/internal/database/repository/user"
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
	// Инициализация конфигурации
	env := config.DetectEnv()
	config.Init(env, "database_init_service", &Cfg)

	// Инициализация логирования
	logger.Init(Cfg.App.LogLevel)

	// Подключение к базе данных
	dsn := Cfg.Database.URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	// Инициализация репозитория
	repo := repository.NewUserRepository(db)

	// Инициализация UseCase
	uc := usecase.NewInitUseCase(repo)

	// Инициализация обработчика
	handler := handler.NewHandler(uc)

	// Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	pb.RegisterDatabaseInitServiceServer(grpcServer, handler)

	logger.Info("Starting database init service on port " + Cfg.App.Port)

	// Запуск gRPC сервера
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
