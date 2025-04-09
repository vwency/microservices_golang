package handler_user_service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
)

type AddUserHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase *usecase.InitUseCase
}

func NewAddUserHandler(uc *usecase.InitUseCase) *AddUserHandler {
	return &AddUserHandler{
		usecase: uc,
	}
}

func (h *AddUserHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	err := h.usecase.AddUser(req.Username, req.Password, req.HashedRt, req.AccessRt, req.Email)
	if err != nil {
		return &pb.AddUserResponse{
			Success: false,
			Message: "User creation failed: " + err.Error(),
		}, nil
	}

	return &pb.AddUserResponse{
		Success: true,
		Message: "User added successfully",
	}, nil
}
