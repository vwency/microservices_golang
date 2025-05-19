package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func decodeInitDatabaseRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.InitRequest)
	return endpoints.InitDatabaseRequest{
		ConfigPath: req.GetConfigPath(),
	}, nil
}

func encodeInitDatabaseResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.InitDatabaseResponse)
	return &pb.InitResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

func (s *grpcServer) InitDatabase(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	_, resp, err := s.initDatabase.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.InitResponse), nil
}
