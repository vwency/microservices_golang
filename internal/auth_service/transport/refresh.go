package transport

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServerRefresh struct {
	refresh grpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func NewGRPCServerRefresh(refreshEndpoint kitendpoint.Endpoint) authv1.AuthServiceServer {
	return &grpcServerRefresh{
		refresh: grpc.NewServer(
			refreshEndpoint,
			decodeRefreshRequest,
			encodeRefreshResponse,
		),
	}
}

func (s *grpcServerRefresh) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	_, resp, err := s.refresh.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh failed: %v", err)
	}
	return resp.(*authv1.RefreshResponse), nil
}

func decodeRefreshRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authv1.RefreshRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid request type")
	}
	return req, nil
}

func encodeRefreshResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*authv1.RefreshResponse)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid response type")
	}
	return resp, nil
}
