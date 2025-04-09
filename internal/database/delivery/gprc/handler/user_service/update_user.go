package handler_user_service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
)

type UpdateUserHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase *usecase.InitUseCase
}

func NewUpdateUserHandler(uc *usecase.InitUseCase) *UpdateUserHandler {
	return &UpdateUserHandler{
		usecase: uc,
	}
}

func (h *UpdateUserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := h.usecase.UpdateUserTokens(req.Username, req.HashedRt, req.AccessRt)
	if err != nil {
		return &pb.UpdateUserResponse{
			Success: false,
			Message: "Token update failed: " + err.Error(),
		}, nil
	}

	return &pb.UpdateUserResponse{
		Success: true,
		Message: "Tokens updated successfully",
	}, nil
}
