package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/hello_service/service"
)

type Endpoints struct {
	SayHelloEndpoint endpoint.Endpoint
}

func MakeEndpoints(s service.HelloService) Endpoints {
	return Endpoints{
		SayHelloEndpoint: makeSayHelloEndpoint(s),
	}
}

func makeSayHelloEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SayHelloRequest)
		resp, err := s.SayHello(ctx, req.Text)
		return SayHelloResponse{Text: resp, Err: err}, nil
	}
}

type SayHelloRequest struct {
	Text string
}

type SayHelloResponse struct {
	Text string
	Err  error
}
