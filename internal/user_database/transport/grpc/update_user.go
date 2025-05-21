package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func decodeUpdateUserRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.UpdateUserRequest)
	return endpoints.UpdateUserRequest{
		UserID:             r.GetUserId(),
		HashedRefreshToken: r.GetHashedRefreshToken(),
		HashedAccessToken:  r.GetHashedAccessToken(),
	}, nil
}

func encodeUpdateUserResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(endpoints.UpdateUserResponse)
	return &pb.UpdateUserResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	_, res, err := s.updateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, GRPCErrorWrapper(err)
	}
	return res.(*pb.UpdateUserResponse), nil
}
