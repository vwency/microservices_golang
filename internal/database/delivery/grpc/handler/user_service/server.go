package handler_user_service_grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server struct that holds all handlers
type Server struct {
	pb.UnimplementedDatabaseInitServiceServer
	addUserHandler     *AddUserHandler
	getUserHandler     *GetUserHandler
	updateUserHandler  *UpdateUserHandler
	getUserByIDHandler *GetUserByIDHandler // Add GetUserByIDHandler here
}

// NewServer initializes the Server with all handlers, including GetUserByIDHandler
func NewServer(uc *user_usecase.UserUsecase, logger *zap.Logger) *Server {
	return &Server{
		addUserHandler:     NewAddUserHandler(uc, logger),
		getUserHandler:     NewGetUserHandler(uc, logger),
		updateUserHandler:  NewUpdateUserHandler(uc, logger),
		getUserByIDHandler: NewGetUserByIDHandler(uc, logger), // Initialize GetUserByIDHandler
	}
}

// Register registers all service methods with the gRPC server
func (s *Server) Register(grpcServer *grpc.Server) {
	pb.RegisterDatabaseInitServiceServer(grpcServer, s)
}

// AddUser handles AddUser requests
func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	return s.addUserHandler.AddUser(ctx, req)
}

// GetUser handles GetUser requests
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.getUserHandler.GetUser(ctx, req)
}

// UpdateUser handles UpdateUser requests
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.updateUserHandler.UpdateUser(ctx, req)
}

// GetUserByID handles GetUserByID requests
func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	return s.getUserByIDHandler.GetUserByID(ctx, req) // Call the handler method
}
