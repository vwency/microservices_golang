package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddUserRequest struct {
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
	UserID             string
}

type AddUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *userService) AddUser(ctx context.Context, request AddUserRequest) (AddUserResponse, error) {
	logger := log.With(s.logger, "method", "AddUser")

	if request.UserID == "" {
		level.Error(logger).Log("msg", "user_id is required")
		return AddUserResponse{}, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		level.Warn(logger).Log("msg", "invalid user_id format", "user_id", request.UserID, "err", err)
		return AddUserResponse{}, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	existingUserByID, err := s.repo.UserRepo.GetUserByID(request.UserID)
	if err != nil && status.Code(err) != codes.NotFound {
		level.Error(logger).Log("msg", "failed to check user existence by ID", "user_id", request.UserID, "err", err)
		return AddUserResponse{}, status.Errorf(codes.Internal, "check user existence by ID failed: %v", err)
	}
	if existingUserByID != nil {
		level.Warn(logger).Log("msg", "user with this ID already exists", "user_id", request.UserID)
		return AddUserResponse{}, status.Error(codes.AlreadyExists, "user already exists")
	}

	existingUser, err := s.repo.UserRepo.GetUserByUsernameOrEmail(request.Username, request.Email)
	if err != nil && status.Code(err) != codes.NotFound {
		level.Error(logger).Log("msg", "failed to check user existence", "username", request.Username, "err", err)
		return AddUserResponse{}, status.Errorf(codes.Internal, "check user existence failed: %v", err)
	}
	if existingUser != nil {
		level.Warn(logger).Log("msg", "user already exists", "username", request.Username)
		return AddUserResponse{}, status.Error(codes.AlreadyExists, "user already exists")
	}

	user := models.User{
		ID:                 userID,
		Username:           request.Username,
		HashedPassword:     request.HashedPassword,
		HashedRefreshToken: request.HashedRefreshToken,
		HashedAccessToken:  request.HashedAccessToken,
	}

	if request.Email != "" {
		user.Email = &request.Email
	}

	if err := s.repo.UserRepo.AddUser(&user); err != nil {
		level.Error(logger).Log("msg", "failed to create user", "username", request.Username, "user_id", request.UserID, "err", err)
		if status.Code(err) != codes.Unknown {
			return AddUserResponse{}, err
		}
		return AddUserResponse{}, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	level.Info(logger).Log("msg", "user created successfully", "username", request.Username, "user_id", request.UserID)
	return AddUserResponse{
		Success: true,
		Message: "user created successfully",
	}, nil
}
