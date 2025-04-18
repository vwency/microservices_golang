package handler_user_service_grpc

import (
	"context"
	"errors"
	"time"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

func NewAddUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *AddUserHandler {
	return &AddUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "add_user")),
	}
}

func (h *AddUserHandler) AddUser(ctx context.Context, req *databasev1.AddUserRequest) (*databasev1.AddUserResponse, error) {
	if req.GetUsername() == "" || req.GetHashedPassword() == "" || req.GetEmail() == "" {
		h.logger.Warn("missing required fields",
			zap.String("username", req.GetUsername()),
			zap.String("email", req.GetEmail()))
		return nil, status.Error(codes.InvalidArgument, "username, hashed password, and email are required")
	}

	params := user_usecase.CreateUserParams{
		UserParams: user_usecase.UserParams{
			Username: req.GetUsername(),
			Email:    req.GetEmail(),
		},
		HashedPassword: req.GetHashedPassword(),
		HashedRt:       req.GetHashedRefreshToken(), // Changed from HashedRt to match proto
		HashedAt:       req.GetHashedAccessToken(),  // Changed from time.Now() to match proto
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	if err := h.uc.CreateUser(params); err != nil {
		h.logger.Error("failed to create user",
			zap.String("username", req.GetUsername()),
			zap.Error(err))

		if errors.Is(err, user_usecase.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	h.logger.Info("user created successfully",
		zap.String("username", req.GetUsername()))

	return &databasev1.AddUserResponse{
		Success: true,
		Message: "User created successfully",
	}, nil
}
