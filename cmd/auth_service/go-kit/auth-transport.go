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

	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      2 * time.Minute,
		MaxConnectionAgeGrace: 15 * time.Second,
		Time:                  10 * time.Second,
		Timeout:               3 * time.Second,
	}
	kaEnforcement := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	opts = append(opts,
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaEnforcement),
		grpc.MaxConcurrentStreams(1000),
	)

	server := grpc.NewServer(opts...)
	transport.RegisterGRPCServer(server, endpoints)
	return server
}
