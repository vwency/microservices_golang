package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

type Endpoints struct {
	AddUser      endpoint.Endpoint
	GetUser      endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
	GetUserByID  endpoint.Endpoint
	InitDatabase endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AddUser:      MakeAddUserEndpoint(s),
		GetUser:      MakeGetUserEndpoint(s),
		UpdateUser:   MakeUpdateUserEndpoint(s),
		DeleteUser:   MakeDeleteUserEndpoint(s),
		GetUserByID:  MakeGetUserByIDEndpoint(s),
		InitDatabase: MakeInitDatabaseEndpoint(s),
	}
}

func MakeGetUserByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GetUserByIDRequest)
		return s.GetUserByID(ctx, req)
	}
}
