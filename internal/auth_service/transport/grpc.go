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
		login:    makeLoginHandler(eps),
		logout:   makeLogoutHandler(eps),
		register: makeRegisterHandler(eps),
		refresh:  makeRefreshHandler(eps),
		validate: makeValidateHandler(eps),
	})
}
