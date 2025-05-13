package transport

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeInitDatabaseHandler(ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) *kitgrpc.Server {
	return kitgrpc.NewServer(
		ep.InitDatabase,
		decodeInitDatabaseRequest,
		encodeInitDatabaseResponse,
		opts...,
	)
}

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
