package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteUserRequest struct {
	UserID string
}

type DeleteUserResponse struct {
	Success bool
}

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	if req.UserID == "" {
		level.Error(logger).Log("msg", "user_id is required")
		return DeleteUserResponse{Success: false}, NewInvalidArgumentError("user_id is required", nil)
	}

	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			level.Warn(logger).Log("msg", "user not found", "userID", req.UserID)
			return DeleteUserResponse{Success: false}, NewNotFoundError("user not found", nil)
		}
		level.Error(logger).Log("msg", "failed to get user", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, NewInternalError("failed to get user", err)
	}

	err = s.repo.UserRepo.DeleteUser(req.UserID)
	if err != nil {
		level.Error(logger).Log("msg", "failed to delete user", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, NewInternalError("failed to delete user", err)
	}

	level.Info(logger).Log("msg", "user deleted successfully", "userID", req.UserID)
	return DeleteUserResponse{
		Success: true,
	}, nil
}
