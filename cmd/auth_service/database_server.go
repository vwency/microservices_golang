package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/vwency/microservices_golang/pkg/config"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

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

	// Загрузка клиентского сертификата
	clientCert, err := tls.LoadX509KeyPair("tls/db_server.crt", "tls/db_server.key")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить клиентский сертификат и ключ: %w", err)
	}

	// Настройка TLS конфигурации
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		ServerName:   "localhost", // Имя сервера должно соответствовать CN в сертификате
	}

	return credentials.NewTLS(config), nil
}

func newDatabaseConnection(cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Загружаем TLS сертификаты
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки TLS сертификатов: %w", err)
	}

	// Используем TLS креды вместо insecure
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
	}

	// Используем grpc.NewClient вместо устаревшего grpc.Dial
	conn, err := grpc.NewClient(cfg.UserDatabase.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к user_database: %w", err)
	}

	// Устанавливаем соединение
	conn.Connect()

	// Ожидаем установления соединения
	state := conn.GetState()
	if state != connectivity.Ready {
		if !conn.WaitForStateChange(ctx, state) {
			conn.Close()
			return nil, fmt.Errorf("connection timeout: context deadline exceeded")
		}
	}

	return conn, nil
}

func newDatabaseClient(conn *grpc.ClientConn) databasev1.DatabaseInitServiceClient {
	return databasev1.NewDatabaseInitServiceClient(conn)
}
