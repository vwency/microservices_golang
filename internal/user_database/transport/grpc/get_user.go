package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	res, err := s.ep.GetUser.Handle(ctx, service.GetUserRequest{
		UserID:   req.GetUserId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Found:              res.Found,
		UserId:             res.UserID,
		Username:           res.Username,
		Email:              res.Email,
		HashedPassword:     res.HashedPassword,
		HashedRefreshToken: res.HashedRefreshToken,
		HashedAccessToken:  res.HashedAccessToken,
	}, nil
}
