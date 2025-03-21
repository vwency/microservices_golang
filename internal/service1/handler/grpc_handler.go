package handler

import (
	"context"

	"github.com/vwency/microservices_golang/proto/service1"
)

type HelloHandler struct {
	service1.UnimplementedHelloServiceServer
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

func (h *HelloHandler) SayHello(ctx context.Context, req *service1.HelloRequest) (*service1.HelloResponse, error) {
	return &service1.HelloResponse{Text: "hello"}, nil
}
