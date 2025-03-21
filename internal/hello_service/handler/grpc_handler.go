package handler

import (
	"context"
    "strings"
	"github.com/vwency/microservices_golang/proto/hello_service"
)

type HelloHandler struct {
	hello_service.UnimplementedHelloServiceServer
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}



func (h *HelloHandler) SayHello(ctx context.Context, req *hello_service.HelloRequest) (*hello_service.HelloResponse, error) {
    if strings.Contains(strings.ToLower(req.GetText()), "hello") {
        return &hello_service.HelloResponse{Text: "hello"}, nil
    }
    return &hello_service.HelloResponse{Text: "None"}, nil
}
