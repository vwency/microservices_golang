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
		return AddUserResponse{}, NewInvalidArgumentError("user_id is required", nil)
	}
	if request.Username == "" {
		level.Error(logger).Log("msg", "username is required")
		return AddUserResponse{}, NewInvalidArgumentError("username is required", nil)
	}
	if request.HashedPassword == "" {
		level.Error(logger).Log("msg", "hashed_password is required")
		return AddUserResponse{}, NewInvalidArgumentError("hashed_password is required", nil)
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		level.Warn(logger).Log("msg", "invalid user_id format", "user_id", request.UserID, "err", err)
		return AddUserResponse{}, NewInvalidArgumentError("invalid user_id format", err)
	}

	existingUserByID, err := s.repo.UserRepo.GetUserByID(request.UserID)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			level.Error(logger).Log("msg", "failed to check user existence by ID", "user_id", request.UserID, "err", err)
			return AddUserResponse{}, NewInternalError("failed to check user existence by ID", err)
		}
	} else if existingUserByID != nil {
		level.Warn(logger).Log("msg", "user with this ID already exists", "user_id", request.UserID)
		return AddUserResponse{}, NewAlreadyExistsError("user with this ID already exists", nil)
	}

	existingUser, err := s.repo.UserRepo.GetUserByUsernameOrEmail(request.Username, request.Email)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			level.Error(logger).Log("msg", "failed to check user existence", "username", request.Username, "err", err)
			return AddUserResponse{}, NewInternalError("failed to check user existence", err)
		}
	} else if existingUser != nil {
		level.Warn(logger).Log("msg", "user already exists", "username", request.Username)
		return AddUserResponse{}, NewAlreadyExistsError("user with this username/email already exists", nil)
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
		return AddUserResponse{}, NewInternalError("failed to create user", err)
	}

	level.Info(logger).Log("msg", "user created successfully", "username", request.Username, "user_id", request.UserID)
	return AddUserResponse{
		Success: true,
		Message: "user created successfully",
	}, nil
}
