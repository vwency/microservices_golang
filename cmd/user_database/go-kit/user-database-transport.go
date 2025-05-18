package gokit

import (
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCServer(creds credentials.TransportCredentials, endpoints endpoints.Endpoints) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	transport.RegisterGRPCServer(grpcServer, endpoints)
	return grpcServer
}
