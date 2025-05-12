package transport

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeDeleteUserHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.DeleteUser,
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
		opts...,
	)
}

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DeleteUserRequest)
	return endpoints.DeleteUserRequest{
		UserID: req.GetUserId(),
	}, nil
}

func encodeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.DeleteUserResponse)
	return &pb.DeleteUserResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}
