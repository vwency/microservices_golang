package transport

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServerLogout struct {
	logout grpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func NewGRPCServerLogout(logoutEndpoint kitendpoint.Endpoint) authv1.AuthServiceServer {
	return &grpcServerLogout{
		logout: grpc.NewServer(
			logoutEndpoint,
			decodeLogoutRequest,
			encodeLogoutResponse,
		),
	}
}

func (s *grpcServerLogout) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	_, resp, err := s.logout.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}
	return resp.(*authv1.LogoutResponse), nil
}

func decodeLogoutRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*authv1.LogoutRequest)
	return req, nil
}

func encodeLogoutResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*authv1.LogoutResponse)
	return resp, nil
}
