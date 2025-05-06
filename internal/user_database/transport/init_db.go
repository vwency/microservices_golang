package transport

import (
	"context"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeInitDatabaseHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.InitDB,
		decodeInitDatabaseRequest,
		encodeInitDatabaseResponse,
		opts...,
	)
}

func decodeInitDatabaseRequest(_ context.Context, req interface{}) (interface{}, error) {
	// No parameters needed for init request
	return service.InitDatabaseRequest{}, nil
}

func encodeInitDatabaseResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(service.InitDatabaseResponse)
	return &pb.InitDatabaseResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}

func (s *grpcServer) InitDB(ctx context.Context, req *pb.InitDatabaseRequest) (*pb.InitDatabaseResponse, error) {
	_, res, err := s.initDB.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.InitDatabaseResponse), nil
}
