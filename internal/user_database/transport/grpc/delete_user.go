package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res, err := s.ep.DeleteUser.Handle(ctx, service.DeleteUserRequest{UserID: req.GetUserId()})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Success: res.Success, Message: res.Message}, nil
}
