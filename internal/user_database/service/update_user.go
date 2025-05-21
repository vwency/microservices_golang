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

	// Проверка существования пользователя
	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				level.Warn(logger).Log("msg", "user not found", "user_id", req.UserID)
				return UpdateUserResponse{}, NewNotFoundError("user not found", err)
			case codes.InvalidArgument:
				level.Error(logger).Log("msg", "invalid argument", "err", err)
				return UpdateUserResponse{}, NewInvalidArgumentError("invalid user ID", err)
			default:
				level.Error(logger).Log("msg", "unexpected error getting user", "err", err)
				return UpdateUserResponse{}, NewInternalError("unexpected error getting user", err)
			}
		}

		level.Error(logger).Log("msg", "unknown error getting user", "err", err)
		return UpdateUserResponse{}, NewInternalError("failed to verify user existence", err)
	}

	// Обновление токенов
	err = s.repo.UserRepo.UpdateUserTokens(req.UserID, req.HashedRefreshToken, req.HashedAccessToken)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				level.Warn(logger).Log("msg", "invalid token format", "err", err)
				return UpdateUserResponse{}, NewInvalidArgumentError("invalid token format", err)
			case codes.Aborted:
				level.Warn(logger).Log("msg", "update aborted", "err", err)
				return UpdateUserResponse{}, NewAbortedError("update aborted", err)
			default:
				level.Error(logger).Log("msg", "failed to update tokens", "err", err)
				return UpdateUserResponse{}, NewInternalError("failed to update tokens", err)
			}
		}

		level.Error(logger).Log("msg", "unknown error updating tokens", "err", err)
		return UpdateUserResponse{}, NewInternalError("unknown error updating tokens", err)
	}

	level.Info(logger).Log("msg", "tokens updated successfully", "user_id", req.UserID)

	return UpdateUserResponse{
		Success: true,
		Message: "tokens updated successfully",
	}, nil
}
