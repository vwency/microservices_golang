package grpc

import (
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
)

func makeAddUserHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.AddUser,
		decodeAddUserRequest,
		encodeAddUserResponse,
		opts...,
	)
}

func makeDeleteUserHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.DeleteUser,
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
		opts...,
	)
}

func makeGetUserByIDHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.GetUserByID,
		decodeGetUserByIDRequest,
		encodeGetUserByIDResponse,
		opts...,
	)
}

func makeInitDatabaseHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.InitDatabase,
		decodeInitDatabaseRequest,
		encodeInitDatabaseResponse,
		opts...,
	)
}

func makeGetUserHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.GetUser,
		decodeGetUserRequest,
		encodeGetUserResponse,
		opts...,
	)
}

func makeUpdateUserHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.UpdateUser,
		decodeUpdateUserRequest,
		encodeUpdateUserResponse,
		opts...,
	)
}
