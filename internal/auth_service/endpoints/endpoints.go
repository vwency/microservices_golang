package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
)

type Endpoints struct {
	RegisterEndpoint endpoint.Endpoint
}

func MakeEndpoints(s service.AuthService) Endpoints {
	return Endpoints{
		RegisterEndpoint: MakeRegisterEndpoint(s),
	}
}
