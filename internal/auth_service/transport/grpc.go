package transport

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	login    kitgrpc.Handler
	logout   kitgrpc.Handler
	register kitgrpc.Handler
	refresh  kitgrpc.Handler
	validate kitgrpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func RegisterGRPCServer(server *grpc.Server, endpoints Endpoints) {
	authv1.RegisterAuthServiceServer(server, &grpcServer{
		login:    kitgrpc.NewServer(endpoints.Login, decodeLoginRequest, encodeLoginResponse),
		logout:   kitgrpc.NewServer(endpoints.Logout, decodeLogoutRequest, encodeLogoutResponse),
		register: kitgrpc.NewServer(endpoints.Register, decodeRegisterRequest, encodeRegisterResponse),
		refresh:  kitgrpc.NewServer(endpoints.Refresh, decodeRefreshRequest, encodeRefreshResponse),
		validate: kitgrpc.NewServer(endpoints.Validate, decodeValidateRequest, encodeValidateResponse),
	})
}

// Методы gRPC сервиса

func (s *grpcServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}
	return resp.(*authv1.LoginResponse), nil
}

func (s *grpcServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	_, resp, err := s.logout.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}
	return resp.(*authv1.LogoutResponse), nil
}

func (s *grpcServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return resp.(*authv1.RegisterResponse), nil
}

func (s *grpcServer) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	_, resp, err := s.refresh.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh failed: %v", err)
	}
	return resp.(*authv1.RefreshResponse), nil
}

func (s *grpcServer) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "validate failed: %v", err)
	}
	return resp.(*authv1.ValidateResponse), nil
}
