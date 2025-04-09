// server.go
package handler_user_service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
)

type Server struct {
	pb.UnimplementedDatabaseInitServiceServer
	DatabaseInitHandler *DatabaseInitHandler
	GetUserHandler      *GetUserHandler
	AddUserHandler      *AddUserHandler
	UpdateUserHandler   *UpdateUserHandler
}

func NewServer(uc *usecase.InitUseCase) *Server {
	return &Server{
		DatabaseInitHandler: NewDatabaseInitHandler(uc),
		GetUserHandler:      NewGetUserHandler(uc),
		AddUserHandler:      NewAddUserHandler(uc),
		UpdateUserHandler:   NewUpdateUserHandler(uc),
	}
}

// Делегируем InitDatabase
func (s *Server) InitDatabase(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	return s.DatabaseInitHandler.InitDatabase(ctx, req)
}

// Делегируем GetUser
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.GetUserHandler.GetUser(ctx, req)
}

// Делегируем AddUser
func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	return s.AddUserHandler.AddUser(ctx, req)
}

// Делегируем UpdateUser
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.UpdateUserHandler.UpdateUser(ctx, req)
}
