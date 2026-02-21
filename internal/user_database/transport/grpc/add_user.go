package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	res, err := s.ep.AddUser.Handle(ctx, service.AddUserRequest{
		Username:           req.GetUsername(),
		Email:              req.GetEmail(),
		HashedPassword:     req.GetHashedPassword(),
		HashedRefreshToken: req.GetHashedRefreshToken(),
		HashedAccessToken:  req.GetHashedAccessToken(),
		UserID:             req.GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddUserResponse{Success: res.Success, Message: res.Message}, nil
}
