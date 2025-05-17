package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadServerTLSCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("tls/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать корневой CA сертификат: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("не удалось добавить CA сертификат в пул")
	}

	serverCert, err := tls.LoadX509KeyPair("tls/auth_client.crt", "tls/auth_client.key")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить серверный сертификат и ключ: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert, // tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func newGRPCServer(endpoints endpoints.Endpoints) *grpc.Server {
	tlsCredentials, err := loadServerTLSCredentials()
	if err != nil {
		fmt.Printf("Предупреждение: не удалось загрузить TLS сертификаты: %v\nСервер будет запущен без TLS!\n", err)
		server := grpc.NewServer()
		transport.RegisterGRPCServer(server, endpoints)
		return server
	}

	server := grpc.NewServer(grpc.Creds(tlsCredentials))
	transport.RegisterGRPCServer(server, endpoints)
	return server
}
