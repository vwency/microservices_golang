package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal error")
)

type DeleteUserRequest struct {
	UserID string
}

type DeleteUserResponse struct {
	Success bool
}

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	if req.UserID == "" {
		return DeleteUserResponse{Success: false}, status.Errorf(codes.InvalidArgument, "userID is required")
	}

	user, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to get user", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if user == nil {
		level.Warn(s.logger).Log("msg", "user not found", "userID", req.UserID)
		return DeleteUserResponse{Success: false}, status.Errorf(codes.NotFound, "user with ID %s not found", req.UserID)
	}

	err = s.repo.UserRepo.DeleteUser(req.UserID)
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to delete user", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	level.Info(s.logger).Log("msg", "user deleted", "userID", req.UserID)
	return DeleteUserResponse{Success: true}, nil
}
