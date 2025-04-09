package handler_user_service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase" // Исправленный импорт
	pb "github.com/vwency/microservices_golang/proto/database"
)

type GetUserHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase *user_usecase.InitUseCase // Используем правильный пакет
}

func NewGetUserHandler(uc *user_usecase.InitUseCase) *GetUserHandler {
	return &GetUserHandler{
		usecase: uc,
	}
}

func (h *GetUserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.usecase.GetUser(req.Username, req.Email)
	if err != nil {
		return &pb.GetUserResponse{
			Found:   false,
			Message: "Error fetching user: " + err.Error(),
		}, nil
	}

	if user == nil {
		// Пользователь не найден
		return &pb.GetUserResponse{
			Found:   false,
			Message: "User not found",
		}, nil
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	return &pb.GetUserResponse{
		Found:    true,
		Username: user.Username,
		Email:    email,
		HashedRt: user.HashedRt,
		Password: user.Password,
		HashedAt: user.HashedAt,
	}, nil
}
