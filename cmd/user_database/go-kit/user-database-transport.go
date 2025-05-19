package gokit

import (
	"time"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	grpcTransport "github.com/vwency/microservices_golang/internal/user_database/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func NewGRPCServer(creds credentials.TransportCredentials, endpoints endpoints.Endpoints) *grpc.Server {
	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second,
		MaxConnectionAge:      2 * time.Minute,
		MaxConnectionAgeGrace: 30 * time.Second,
		Time:                  30 * time.Second,
		Timeout:               10 * time.Second,
	}

	keepalivePolicy := keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second,
		PermitWithoutStream: true,
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.KeepaliveParams(keepaliveParams),
		grpc.KeepaliveEnforcementPolicy(keepalivePolicy),
	)

	grpcTransport.RegisterGRPCServer(grpcServer, endpoints)
	return grpcServer
}
