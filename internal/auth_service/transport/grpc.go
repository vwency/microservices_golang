package transport

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc"
)

type grpcServer struct {
	login    kitgrpc.Handler
	logout   kitgrpc.Handler
	register kitgrpc.Handler
	refresh  kitgrpc.Handler
	validate kitgrpc.Handler
	authv1.UnimplementedAuthServiceServer
}

func RegisterGRPCServer(server *grpc.Server, eps endpoints.Endpoints) {
	authv1.RegisterAuthServiceServer(server, &grpcServer{
		login:    kitgrpc.NewServer(eps.Login, decodeLoginRequest, encodeLoginResponse),
		logout:   kitgrpc.NewServer(eps.Logout, decodeLogoutRequest, encodeLogoutResponse),
		register: kitgrpc.NewServer(eps.Register, decodeRegisterRequest, encodeRegisterResponse),
		refresh:  kitgrpc.NewServer(eps.Refresh, decodeRefreshRequest, encodeRefreshResponse),
		validate: kitgrpc.NewServer(eps.ValidateAccessToken, decodeValidateRequest, encodeValidateResponse),
	})
}
