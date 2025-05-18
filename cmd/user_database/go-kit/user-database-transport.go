package gokit

import (
	"time"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func NewGRPCServer(creds credentials.TransportCredentials, endpoints endpoints.Endpoints) *grpc.Server {
	// Настройки keepalive для сервера
	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second, // Максимальное время бездействия соединения
		MaxConnectionAge:      2 * time.Minute,  // Максимальное время жизни соединения
		MaxConnectionAgeGrace: 30 * time.Second, // Дополнительное время для завершения операций
		Time:                  30 * time.Second, // Период отправки ping
		Timeout:               10 * time.Second, // Таймаут ожидания ответа на ping
	}

	// Настройки принудительного keepalive
	keepalivePolicy := keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second, // Минимальное время между ping от клиента
		PermitWithoutStream: true,             // Разрешить ping даже когда нет потоков
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.KeepaliveParams(keepaliveParams),
		grpc.KeepaliveEnforcementPolicy(keepalivePolicy),
	)

	transport.RegisterGRPCServer(grpcServer, endpoints)
	return grpcServer
}
