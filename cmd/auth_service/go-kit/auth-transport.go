package gokit

import (
	"fmt"

	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCServer(
	endpoints endpoints.Endpoints,
	tlsCredentials credentials.TransportCredentials,
) *grpc.Server {
	if tlsCredentials == nil {
		fmt.Println("Warning: Running server without TLS")
		server := grpc.NewServer()
		transport.RegisterGRPCServer(server, endpoints)
		return server
	}

	server := grpc.NewServer(grpc.Creds(tlsCredentials))
	transport.RegisterGRPCServer(server, endpoints)
	return server
}
