package transport

import (
	"context"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	pb "github.com/vwency/microservices_golang/proto/user_database"
)

func makeGetUserByIDHandler(ep endpoints.Endpoints, opts ...gokitgrpc.ServerOption) *gokitgrpc.Server {
	return gokitgrpc.NewServer(
		ep.GetUserByID,
		decodeGetUserByIDRequest,
		encodeGetUserByIDResponse,
		opts...,
	)
}

func decodeGetUserByIDRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.GetUserByIDRequest)
	return service.GetUserByIDRequest{
		UserID: r.GetUserId(),
	}, nil
}

func encodeGetUserByIDResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(service.GetUserByIDResponse)
	return &pb.GetUserByIDResponse{
		UserId:             r.UserID,
		Username:           r.Username,
		Email:              r.Email,
		HashedPassword:     r.HashedPassword,
		HashedRefreshToken: r.HashedRefreshToken,
		HashedAccessToken:  r.HashedAccessToken,
	}, nil
}

func (s *grpcServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	_, res, err := s.getUserByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetUserByIDResponse), nil
}
