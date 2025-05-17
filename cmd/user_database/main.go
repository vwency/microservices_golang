package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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

	// Load TLS credentials
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		logger.Fatal("Failed to load TLS credentials", zap.Error(err))
	}

	// Create gRPC server with TLS
	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	transport.RegisterGRPCServer(grpcServer, eps)

	lis, err := net.Listen("tcp", "0.0.0.0:"+Cfg.App.Port)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	logger.Info("Starting server",
		zap.String("service", Cfg.App.ServiceName),
		zap.String("port", Cfg.App.Port),
		zap.Bool("TLS", true),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("gRPC server failed", zap.Error(err))
	}
}

// Функция для загрузки TLS сертификатов и создания конфигурации
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Загрузка корневого CA сертификата
	pemServerCA, err := os.ReadFile("tls/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать корневой CA сертификат: %w", err)
	}

	// Создание пула сертификатов, добавление CA сертификата
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("не удалось добавить CA сертификат в пул")
	}

	// Загрузка серверного сертификата и ключа
	serverCert, err := tls.LoadX509KeyPair("tls/db_server.crt", "tls/db_server.key")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить серверный сертификат и ключ: %w", err)
	}

	// Настройка TLS конфигурации
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
		ServerName:   "localhost", // Должно соответствовать CN в сертификате
	}

	return credentials.NewTLS(config), nil
}
