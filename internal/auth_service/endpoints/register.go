package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
)

func MakeRegisterEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*service.RegisterRequest)
		return s.Register(ctx, req)
	}
}
