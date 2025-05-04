package transport

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServerLogin struct {
	login grpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func NewGRPCServerLogin(loginEndpoint kitendpoint.Endpoint) authv1.AuthServiceServer {
	return &grpcServerLogin{
		login: grpc.NewServer(
			loginEndpoint,
			decodeLoginRequest,
			encodeLoginResponse,
		),
	}
}

func (s *grpcServerLogin) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}
	return resp.(*authv1.LoginResponse), nil
}

func decodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authv1.LoginRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid request type")
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*authv1.LoginResponse)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid response type")
	}
	return resp, nil
}
