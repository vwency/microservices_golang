package transport

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServerRegister struct {
	register grpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func NewGRPCServerRegister(registerEndpoint kitendpoint.Endpoint) authv1.AuthServiceServer {
	return &grpcServerRegister{
		register: grpc.NewServer(
			registerEndpoint,
			decodeRegisterRequest,
			encodeRegisterResponse,
		),
	}
}

func (s *grpcServerRegister) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return resp.(*authv1.RegisterResponse), nil
}

func decodeRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*authv1.RegisterRequest)
	return &service.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}, nil
}

func encodeRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*service.RegisterResponse)
	return &authv1.RegisterResponse{
		UserId:       resp.UserID,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	}, nil
}
