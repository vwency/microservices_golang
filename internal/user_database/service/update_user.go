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
	const successMessage = "tokens updated successfully"

	// Validate required fields
	if req.UserID == "" {
		msg := "userID is required"
		level.Error(logger).Log("msg", msg)
		return UpdateUserResponse{
			Success: false,
			Message: msg,
		}, status.Errorf(codes.InvalidArgument, msg)
	}

	if req.HashedRefreshToken == "" || req.HashedAccessToken == "" {
		msg := "both refresh and access tokens are required"
		level.Error(logger).Log("msg", msg)
		return UpdateUserResponse{
			Success: false,
			Message: msg,
		}, status.Errorf(codes.InvalidArgument, msg)
	}

	// Check if user exists
	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			msg := "user not found"
			level.Warn(logger).Log("msg", msg, "userID", req.UserID)
			return UpdateUserResponse{
				Success: false,
				Message: msg,
			}, status.Errorf(codes.NotFound, msg)
		}

		msg := "failed to verify user existence"
		level.Error(logger).Log("msg", msg, "userID", req.UserID, "err", err)
		return UpdateUserResponse{
			Success: false,
			Message: msg,
		}, status.Errorf(codes.Internal, msg)
	}

	// Update tokens
	if err := s.repo.UserRepo.UpdateUserTokens(
		req.UserID,
		req.HashedRefreshToken,
		req.HashedAccessToken,
	); err != nil {
		msg := "failed to update tokens"
		level.Error(logger).Log("msg", msg, "userID", req.UserID, "err", err)
		return UpdateUserResponse{
			Success: false,
			Message: msg,
		}, status.Errorf(codes.Internal, msg)
	}

	level.Info(logger).Log(
		"msg", successMessage,
		"userID", req.UserID,
	)

	return UpdateUserResponse{
		Success: true,
		Message: successMessage,
	}, nil
}
