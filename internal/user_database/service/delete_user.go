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
	Message string
}

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	if req.UserID == "" {
		err := NewInvalidArgumentError("user_id is required", nil)
		level.Error(logger).Log("msg", err.Error())
		return DeleteUserResponse{Success: false, Message: err.Message}, err
	}

	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			wrappedErr := NewNotFoundError("user not found", nil)
			level.Warn(logger).Log("msg", wrappedErr.Error(), "userID", req.UserID)
			return DeleteUserResponse{Success: false, Message: wrappedErr.Message}, wrappedErr
		}
		wrappedErr := NewInternalError("failed to get user", err)
		level.Error(logger).Log("msg", wrappedErr.Error(), "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false, Message: wrappedErr.Message}, wrappedErr
	}

	err = s.repo.UserRepo.DeleteUser(req.UserID)
	if err != nil {
		wrappedErr := NewInternalError("failed to delete user", err)
		level.Error(logger).Log("msg", wrappedErr.Error(), "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false, Message: wrappedErr.Message}, wrappedErr
	}

	level.Info(logger).Log("msg", "user deleted successfully", "userID", req.UserID)
	return DeleteUserResponse{
		Success: true,
		Message: "user deleted successfully",
	}, nil
}
