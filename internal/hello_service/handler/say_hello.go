package handler_hello

import (
	usecase_hello "github.com/vwency/microservices_golang/internal/hello_service/usecase"
	"github.com/vwency/microservices_golang/proto/hello_service"
)

type HelloHandler struct {
	hello_service.UnimplementedHelloServiceServer
	usecase usecase_hello.HelloUsecase
}

func NewHelloHandler(usecase usecase_hello.HelloUsecase) *HelloHandler {
	return &HelloHandler{
		usecase: usecase,
	}
}
