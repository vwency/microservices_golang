package handler_user_service_grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

func NewGetUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *GetUserHandler {
	return &GetUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "get_user")),
	}
}

func (h *GetUserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.GetUsername() == "" && req.GetEmail() == "" {
		h.logger.Warn("empty request parameters")
		return nil, status.Error(codes.InvalidArgument, "username or email must be provided")
	}

	params := user_usecase.UserParams{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	user, err := h.uc.GetUser(params)
	if err != nil {
		// Проверяем, является ли ошибка NotFound
		if status.Code(err) == codes.NotFound {
			h.logger.Info("user not found",
				zap.String("username", req.GetUsername()),
				zap.String("email", req.GetEmail()))
			return &pb.GetUserResponse{
				Found:   false,
				Message: "User not found",
			}, nil
		}

		h.logger.Error("failed to get user",
			zap.String("username", req.GetUsername()),
			zap.String("email", req.GetEmail()),
			zap.Error(err))
		return nil, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	userID := ""
	if user.ID != [16]byte{} {
		userID = user.ID.String()
	}

	return &pb.GetUserResponse{
		Found:              true,
		UserId:             userID,
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
		Message:            "User found",
	}, nil
}
