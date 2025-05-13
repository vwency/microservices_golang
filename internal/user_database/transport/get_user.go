package transport

import (
	"context"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeGetUserHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.GetUser,
		decodeGetUserRequest,
		encodeGetUserResponse,
		opts...,
	)
}

func decodeGetUserRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.GetUserRequest)
	return endpoints.GetUserRequest{
		UserID:   r.GetUserId(),
		Username: r.GetUsername(),
		Email:    r.GetEmail(),
	}, nil
}

func encodeGetUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(endpoints.GetUserResponse)
	if r.Err != nil {
		return nil, r.Err
	}
	return &pb.GetUserResponse{
		Found:              r.Found,
		UserId:             r.UserID,
		Username:           r.Username,
		Email:              r.Email,
		HashedPassword:     r.HashedPassword,
		HashedRefreshToken: r.HashedRefreshToken,
		HashedAccessToken:  r.HashedAccessToken,
	}, nil
}

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, res, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetUserResponse), nil
}
