package handler

import (
	"context"

	pb "github.com/vwency/microservices_golang/proto/database"

	"github.com/vwency/microservices_golang/internal/database/usecase"
)

type DatabaseInitHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase usecase.DatabaseInit
}

func NewHandler(uc usecase.DatabaseInit) *DatabaseInitHandler {
	return &DatabaseInitHandler{
		usecase: uc,
	}
}

func (h *DatabaseInitHandler) InitDatabase(ctx context.Context, req *pb.InitRequest) (*pb.InitResponse, error) {
	err := h.usecase.InitDatabase()
	if err != nil {
		return &pb.InitResponse{
			Success: false,
			Message: "Initialization failed: " + err.Error(),
		}, nil
	}

	return &pb.InitResponse{
		Success: true,
		Message: "Database initialized successfully",
	}, nil
}

func (h *DatabaseInitHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.usecase.GetUser(req.GetUsername(), req.GetEmail())
	if err != nil || user == nil {
		return &pb.GetUserResponse{Found: false}, nil
	}

	return &pb.GetUserResponse{
		Found:    true,
		Username: user.Username,
		Email:    user.Email,
		HashedRt: user.HashedRt,
		Password: user.Password,
		HashedAt: user.HashedAt,
	}, nil
}

func (h *DatabaseInitHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	err := h.usecase.AddUser(req.GetUsername(), req.GetPassword(), req.GetHashedRt(), req.GetAccessRt())
	if err != nil {
		return &pb.AddUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AddUserResponse{
		Success: true,
		Message: "User added successfully",
	}, nil
}

func (h *DatabaseInitHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := h.usecase.UpdateUserTokens(req.GetUsername(), req.GetHashedRt(), req.GetAccessRt())
	if err != nil {
		return &pb.UpdateUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateUserResponse{
		Success: true,
		Message: "User tokens updated successfully",
	}, nil
}
