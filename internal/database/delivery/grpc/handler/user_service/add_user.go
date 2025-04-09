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

func (h *AddUserHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	// Check if required fields are missing
	if req.GetUsername() == "" || req.GetHashedPassword() == "" || req.GetEmail() == "" {
		h.logger.Warn("missing required fields",
			zap.String("username", req.GetUsername()),
			zap.String("email", req.GetEmail()))
		return nil, status.Error(codes.InvalidArgument, "username, hashed password, and email are required")
	}

	// Prepare the parameters for creating the user
	params := user_usecase.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: req.GetHashedPassword(), // Changed to GetHashedPassword
		HashedRt:       req.GetHashedRt(),
		HashedAt:       req.GetAccessRt(), // Assuming AccessRt is the hashed at timestamp
		Email:          req.GetEmail(),
	}

	// Attempt to create the user
	if err := h.uc.CreateUser(params); err != nil {
		h.logger.Error("failed to create user",
			zap.String("username", req.GetUsername()),
			zap.Error(err))

		// Check if the user already exists
		if errors.Is(err, user_usecase.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		// For all other errors, return a generic failure message
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// Log successful user creation
	h.logger.Info("user created successfully",
		zap.String("username", req.GetUsername()))

	// Return a successful response
	return &pb.AddUserResponse{
		Success: true,
		Message: "User created successfully",
	}, nil
}
