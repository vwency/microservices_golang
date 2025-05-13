package handler_user_service_grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/user_database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDatabaseInitServiceServer
	addUserHandler     *AddUserHandler
	getUserHandler     *GetUserHandler
	updateUserHandler  *UpdateUserHandler
	getUserByIDHandler *GetUserByIDHandler
}

func NewServer(uc *user_usecase.UserUsecase, logger *zap.Logger) *Server {
	return &Server{
		addUserHandler:     NewAddUserHandler(uc, logger),
		getUserHandler:     NewGetUserHandler(uc, logger),
		updateUserHandler:  NewUpdateUserHandler(uc, logger),
		getUserByIDHandler: NewGetUserByIDHandler(uc, logger),
	}
}

func (s *Server) Register(grpcServer *grpc.Server) {
	pb.RegisterDatabaseInitServiceServer(grpcServer, s)
}

func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	return s.addUserHandler.AddUser(ctx, req)
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.getUserHandler.GetUser(ctx, req)
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.updateUserHandler.UpdateUser(ctx, req)
}

func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	return s.getUserByIDHandler.GetUserByID(ctx, req)
}
