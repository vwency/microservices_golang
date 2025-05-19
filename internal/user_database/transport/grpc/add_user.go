package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func decodeAddUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.AddUserRequest)
	return endpoints.AddUserRequest{
		Username:           req.GetUsername(),
		Email:              req.GetEmail(),
		HashedPassword:     req.GetHashedPassword(),
		HashedRefreshToken: req.GetHashedRefreshToken(),
		HashedAccessToken:  req.GetHashedAccessToken(),
		UserID:             req.GetUserId(),
	}, nil
}

func encodeAddUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.AddUserResponse)
	return &pb.AddUserResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

func (s *grpcServer) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	_, resp, err := s.addUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, ConvertToGRPCError(err)
	}
	return resp.(*pb.AddUserResponse), nil
}
