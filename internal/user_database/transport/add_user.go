package transport

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeAddUserHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.AddUser,
		decodeAddUserRequest,
		encodeAddUserResponse,
		opts...,
	)
}

func decodeAddUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.AddUserRequest)
	return endpoints.AddUserRequest{
		Username:           req.GetUsername(),
		Email:              req.GetEmail(),
		HashedPassword:     req.GetHashedPassword(),
		HashedRefreshToken: req.GetHashedRefreshToken(),
		HashedAccessToken:  req.GetHashedAccessToken(),
		UserID:             req.GetUserId(),
	}, nil
}

func encodeAddUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.AddUserResponse)
	return &pb.AddUserResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}
