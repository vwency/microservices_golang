package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
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

func (s *userService) AddUser(ctx context.Context, req AddUserRequest) (AddUserResponse, error) {
	logger := log.With(s.logger, "method", "AddUser")

	if req.UserID == "" {
		level.Error(logger).Log("msg", "user_id is required")
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.InvalidArgument, "user_id is required")
	}

	if req.Username == "" {
		level.Error(logger).Log("msg", "username is required")
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.InvalidArgument, "username is required")
	}

	if req.HashedPassword == "" {
		level.Error(logger).Log("msg", "hashed_password is required")
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.InvalidArgument, "hashed_password is required")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		level.Warn(logger).Log("msg", "invalid user_id format", "user_id", req.UserID, "err", err)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.InvalidArgument, "invalid user_id format: "+err.Error())
	}

	_, err = s.repo.UserRepo.GetUserByID(req.UserID)
	if err == nil {
		level.Warn(logger).Log("msg", "user with this ID already exists", "user_id", req.UserID)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.AlreadyExists, "user with this ID already exists")
	} else if status.Code(err) != codes.NotFound {
		level.Error(logger).Log("msg", "failed to check user existence by ID", "user_id", req.UserID, "err", err)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.Internal, "failed to check user existence by ID: "+err.Error())
	}

	_, err = s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err == nil {
		level.Warn(logger).Log("msg", "user with this username/email already exists", "username", req.Username, "email", req.Email)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.AlreadyExists, "user with this username/email already exists")
	} else if status.Code(err) != codes.NotFound {
		level.Error(logger).Log("msg", "failed to check user existence", "username", req.Username, "email", req.Email, "err", err)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.Internal, "failed to check user existence: "+err.Error())
	}

	user := models.User{
		ID:                 userID,
		Username:           req.Username,
		HashedPassword:     req.HashedPassword,
		HashedRefreshToken: req.HashedRefreshToken,
		HashedAccessToken:  req.HashedAccessToken,
	}
	if req.Email != "" {
		user.Email = &req.Email
	}

	if err := s.repo.UserRepo.AddUser(&user); err != nil {
		level.Error(logger).Log("msg", "failed to create user", "user_id", req.UserID, "username", req.Username, "err", err)
		return AddUserResponse{Success: false}, error_hndl.NewError(codes.Internal, "failed to create user: "+err.Error())
	}

	level.Info(logger).Log("msg", "user created successfully", "user_id", req.UserID, "username", req.Username)
	return AddUserResponse{
		Success: true,
		Message: "user created successfully",
	}, nil
}
