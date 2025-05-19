package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DeleteUserRequest)
	return endpoints.DeleteUserRequest{
		UserID: req.GetUserId(),
	}, nil
}

func encodeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.DeleteUserResponse)
	return &pb.DeleteUserResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, resp, err := s.deleteUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteUserResponse), nil
}
