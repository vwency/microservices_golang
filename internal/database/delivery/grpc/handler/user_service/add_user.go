package handler_user_service

import (
	"context"
	"errors"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddUserHandler обрабатывает запросы на добавление пользователя
type AddUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

// NewAddUserHandler создает новый экземпляр обработчика
func NewAddUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *AddUserHandler {
	return &AddUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "add_user")),
	}
}

func (h *AddUserHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	// Минимальная проверка обязательных полей
	if req.GetUsername() == "" || req.GetPassword() == "" || req.GetEmail() == "" {
		h.logger.Warn("missing required fields",
			zap.String("username", req.GetUsername()),
			zap.String("email", req.GetEmail()))
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}

	params := user_usecase.CreateUserParams{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		HashedRt: req.GetHashedRt(),
		HashedAt: req.GetAccessRt(),
		Email:    req.GetEmail(),
	}

	if err := h.uc.CreateUser(params); err != nil {
		h.logger.Error("failed to create user",
			zap.String("username", req.GetUsername()),
			zap.Error(err))

		if errors.Is(err, user_usecase.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		// Для всех остальных ошибок возвращаем общее сообщение
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	h.logger.Info("user created successfully",
		zap.String("username", req.GetUsername()))

	return &pb.AddUserResponse{
		Success: true,
		Message: "User created successfully",
	}, nil
}
