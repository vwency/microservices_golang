package grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func (s *grpcServer) InitDatabase(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	res, err := s.ep.InitDatabase.Handle(ctx, service.InitDatabaseRequest{ConfigPath: req.GetConfigPath()})
	if err != nil {
		return nil, err
	}
	return &pb.InitResponse{Success: res.Success, Message: "Database initialized successfully"}, nil
}
