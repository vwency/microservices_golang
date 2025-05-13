package transport

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	pb "github.com/vwency/microservices_golang/proto/user_database"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedDatabaseInitServiceServer
	addUser      kitgrpc.Handler
	getUser      kitgrpc.Handler
	getUserByID  kitgrpc.Handler
	updateUser   kitgrpc.Handler
	deleteUser   kitgrpc.Handler
	initDatabase kitgrpc.Handler
}

func RegisterGRPCServer(server *grpc.Server, ep endpoints.Endpoints, opts ...kitgrpc.ServerOption) {
	pb.RegisterDatabaseInitServiceServer(server, &grpcServer{
		addUser:      makeAddUserHandler(ep, opts...),
		getUser:      makeGetUserHandler(ep, opts...),
		getUserByID:  makeGetUserByIDHandler(ep, opts...),
		updateUser:   makeUpdateUserHandler(ep, opts...),
		deleteUser:   makeDeleteUserHandler(ep, opts...),
		initDatabase: makeInitDatabaseHandler(ep, opts...),
	})
}
