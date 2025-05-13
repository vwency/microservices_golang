package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
)

type Endpoints struct {
	Login               endpoint.Endpoint
	Logout              endpoint.Endpoint
	Register            endpoint.Endpoint
	Refresh             endpoint.Endpoint
	ValidateAccessToken endpoint.Endpoint
}

func MakeEndpoints(s service.AuthService) Endpoints {
	return Endpoints{
		Login:               MakeLoginEndpoint(s),
		Logout:              MakeLogoutEndpoint(s),
		Register:            MakeRegisterEndpoint(s),
		Refresh:             MakeRefreshEndpoint(s),
		ValidateAccessToken: MakeValidateEndpoint(s),
	}
}
