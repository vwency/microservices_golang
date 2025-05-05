package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login    endpoint.Endpoint
	Logout   endpoint.Endpoint
	Register endpoint.Endpoint
	Refresh  endpoint.Endpoint
	Validate endpoint.Endpoint
}

func MakeEndpoints(s AuthService) Endpoints {
	return Endpoints{
		Login:    MakeLoginEndpoint(s),
		Logout:   MakeLogoutEndpoint(s),
		Register: MakeRegisterEndpoint(s),
		Refresh:  MakeRefreshEndpoint(s),
		Validate: MakeValidateEndpoint(s),
	}
}
