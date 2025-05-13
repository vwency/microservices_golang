package service

import (
	"context"

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
	if req.UserID == "" {
		return UpdateUserResponse{
			Success: false,
			Message: "userID is required",
		}, status.Error(codes.InvalidArgument, "userID is required")
	}

	if req.HashedRefreshToken == "" || req.HashedAccessToken == "" {
		return UpdateUserResponse{
			Success: false,
			Message: "both tokens are required",
		}, status.Error(codes.InvalidArgument, "both tokens are required")
	}

	if _, err := s.repo.UserRepo.GetUserByID(req.UserID); err != nil {
		level.Error(s.logger).Log(
			"msg", "failed to get user for update",
			"userID", req.UserID,
			"err", err,
		)
		return UpdateUserResponse{
			Success: false,
			Message: "user not found",
		}, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	if err := s.repo.UserRepo.UpdateUserTokens(req.UserID, req.HashedRefreshToken, req.HashedAccessToken); err != nil {
		level.Error(s.logger).Log(
			"msg", "failed to update user tokens",
			"userID", req.UserID,
			"err", err,
		)
		return UpdateUserResponse{
			Success: false,
			Message: "failed to update tokens",
		}, status.Errorf(codes.Internal, "failed to update tokens: %v", err)
	}

	level.Info(s.logger).Log(
		"msg", "user tokens updated successfully",
		"userID", req.UserID,
	)

	return UpdateUserResponse{
		Success: true,
		Message: "tokens updated successfully",
	}, nil
}
