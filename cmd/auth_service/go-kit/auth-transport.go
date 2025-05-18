package gokit

import (
	"fmt"
	"time"

	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func NewGRPCServer(
	endpoints endpoints.Endpoints,
	tlsCredentials credentials.TransportCredentials,
) *grpc.Server {
	var opts []grpc.ServerOption

	if tlsCredentials == nil {
		fmt.Println("Warning: Running server without TLS")
	} else {
		opts = append(opts, grpc.Creds(tlsCredentials))
	}

	// Keepalive параметры для улучшения производительности
	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Minute, // закрывать неактивные соединения через 15 минут
		MaxConnectionAge:      24 * time.Hour,   // максимальный возраст соединения
		MaxConnectionAgeGrace: 5 * time.Minute,  // grace period перед закрытием
		Time:                  10 * time.Minute, // время отправки ping клиенту
		Timeout:               20 * time.Second, // время ожидания pong от клиента
	}

	opts = append(opts, grpc.KeepaliveParams(kaParams))

	server := grpc.NewServer(opts...)
	transport.RegisterGRPCServer(server, endpoints)
	return server
}
