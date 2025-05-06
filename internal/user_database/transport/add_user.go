package transport

import (
	"context"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeAddUserHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.AddUser,
		decodeAddUserRequest,
		encodeAddUserResponse,
		opts...,
	)
}

func decodeAddUserRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.AddUserRequest)
	return service.AddUserRequest{
		UserID:             r.GetUserId(),
		Username:           r.GetUsername(),
		Email:              r.GetEmail(),
		HashedPassword:     r.GetHashedPassword(),
		HashedRefreshToken: r.GetHashedRefreshToken(),
		HashedAccessToken:  r.GetHashedAccessToken(),
	}, nil
}

func encodeAddUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(service.AddUserResponse)
	return &pb.AddUserResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}

func (s *grpcServer) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	_, res, err := s.addUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.AddUserResponse), nil
}
