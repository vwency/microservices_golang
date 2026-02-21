package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	res, err := s.ep.UpdateUser.Handle(ctx, service.UpdateUserRequest{
		UserID:             req.GetUserId(),
		HashedRefreshToken: req.GetHashedRefreshToken(),
		HashedAccessToken:  req.GetHashedAccessToken(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{Success: res.Success, Message: res.Message}, nil
}
