package grpc

import (
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedDatabaseInitServiceServer
	ep endpoints.Endpoints
}

func RegisterGRPCServer(s *grpc.Server, ep endpoints.Endpoints) {
	pb.RegisterDatabaseInitServiceServer(s, &grpcServer{ep: ep})
}
