package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type Endpoints struct {
	AddUser    endpoint.Endpoint
	GetUser    endpoint.Endpoint
	UpdateUser endpoint.Endpoint
	DeleteUser endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AddUser:    MakeAddUserEndpoint(s),
		GetUser:    MakeGetUserEndpoint(s),
		UpdateUser: MakeUpdateUserEndpoint(s),
		DeleteUser: MakeDeleteUserEndpoint(s),
	}
}
