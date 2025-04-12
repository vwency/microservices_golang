package handler_hello

import (
	"context"

	"github.com/vwency/microservices_golang/proto/hello_service"
	"google.golang.org/grpc"
)

func (h *HelloHandler) RegisterService(server *grpc.Server) {
	hello_service.RegisterHelloServiceServer(server, h)
}

func (h *HelloHandler) SayHello(ctx context.Context, req *hello_service.HelloRequest) (*hello_service.HelloResponse, error) {
	responseText := h.usecase.ProcessGreeting(req.GetText())
	return &hello_service.HelloResponse{Text: responseText}, nil
}
