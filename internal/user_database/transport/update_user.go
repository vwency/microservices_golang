package transport

import (
	"context"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeUpdateUserHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.UpdateUser,
		decodeUpdateUserRequest,
		encodeUpdateUserResponse,
		opts...,
	)
}

func decodeUpdateUserRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.UpdateUserRequest)
	return service.UpdateUserRequest{
		UserID:             r.GetUserId(),
		Username:           r.GetUsername(),
		Email:              r.GetEmail(),
		HashedPassword:     r.GetHashedPassword(),
		HashedRefreshToken: r.GetHashedRefreshToken(),
		HashedAccessToken:  r.GetHashedAccessToken(),
	}, nil
}

func encodeUpdateUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(service.UpdateUserResponse)
	return &pb.UpdateUserResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	_, res, err := s.updateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.UpdateUserResponse), nil
}
