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
		return DeleteUserResponse{Success: false}, ErrInvalidArgument
	}

	err := s.repo.UserRepo.DeleteUser(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			level.Warn(s.logger).Log("msg", "user not found", "userID", req.UserID)
			return DeleteUserResponse{Success: false}, ErrNotFound
		}
		level.Error(s.logger).Log("msg", "deletion failed", "userID", req.UserID, "err", err)
		return DeleteUserResponse{Success: false}, ErrInternal
	}

	level.Info(s.logger).Log("msg", "user deleted", "userID", req.UserID)
	return DeleteUserResponse{Success: true}, nil
}
