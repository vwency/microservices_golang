package handler

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
)

type DatabaseInitHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	usecase *usecase.InitUseCase
}

func NewHandler(uc *usecase.InitUseCase) *DatabaseInitHandler {
	return &DatabaseInitHandler{
		usecase: uc,
	}
}

// InitDatabase выполняет инициализацию базы данных
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

// GetUser возвращает пользователя по имени или email
func (h *DatabaseInitHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.usecase.GetUser(req.Username, req.Email)
	if err != nil {
		// Сообщаем, что пользователя не найдено
		return &pb.GetUserResponse{
			Found: false,
		}, nil
	}

	// Возвращаем информацию о пользователе
	return &pb.GetUserResponse{
		Found:    true,
		Username: user.Username,
		Email:    user.Email,
		HashedRt: user.HashedRt,
		Password: user.Password,
		HashedAt: user.HashedAt,
	}, nil
}

// AddUser добавляет нового пользователя
func (h *DatabaseInitHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	err := h.usecase.AddUser(req.Username, req.Password, req.HashedRt, req.AccessRt)
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

// UpdateUser обновляет токены пользователя
func (h *DatabaseInitHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
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
