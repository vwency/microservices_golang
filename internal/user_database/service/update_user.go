package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateUserRequest struct {
	UserID             string
	HashedRefreshToken string
	HashedAccessToken  string
}

type UpdateUserResponse struct {
	Success bool
	Message string
}

func (s *userService) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	logger := log.With(s.logger, "method", "UpdateUser")

	if req.UserID == "" {
		level.Error(logger).Log("msg", "user_id is required")
		return UpdateUserResponse{}, NewInvalidArgumentError("user_id is required", nil)
	}

	if req.HashedRefreshToken == "" || req.HashedAccessToken == "" {
		level.Error(logger).Log("msg", "both refresh and access tokens are required")
		return UpdateUserResponse{}, NewInvalidArgumentError("both refresh and access tokens are required", nil)
	}

	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			level.Warn(logger).Log("msg", "user not found", "user_id", req.UserID)
			return UpdateUserResponse{}, NewNotFoundError("user not found", nil)
		}
		level.Error(logger).Log("msg", "failed to verify user existence", "user_id", req.UserID, "err", err)
		return UpdateUserResponse{}, NewInternalError("failed to verify user existence", err)
	}

	err = s.repo.UserRepo.UpdateUserTokens(req.UserID, req.HashedRefreshToken, req.HashedAccessToken)
	if err != nil {
		level.Error(logger).Log("msg", "failed to update tokens", "user_id", req.UserID, "err", err)
		return UpdateUserResponse{}, NewInternalError("failed to update tokens", err)
	}

	level.Info(logger).Log("msg", "tokens updated successfully", "user_id", req.UserID)

	return UpdateUserResponse{
		Success: true,
		Message: "tokens updated successfully",
	}, nil
}
