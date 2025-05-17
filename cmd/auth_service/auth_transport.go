package main

import (
	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/transport"
	"google.golang.org/grpc"
)

func newGRPCServer(endpoints endpoints.Endpoints) *grpc.Server {
	server := grpc.NewServer()
	transport.RegisterGRPCServer(server, endpoints)
	return server
}
