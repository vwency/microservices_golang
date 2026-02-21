package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	res, err := s.ep.GetUserByID.Handle(ctx, service.GetUserByIDRequest{UserID: req.GetUserId()})
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByIDResponse{
		UserId:             res.UserID,
		Username:           res.Username,
		Email:              res.Email,
		HashedPassword:     res.HashedPassword,
		HashedRefreshToken: res.HashedRefreshToken,
		HashedAccessToken:  res.HashedAccessToken,
	}, nil
}
