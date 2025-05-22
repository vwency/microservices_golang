package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
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
		err := error_hndl.NewError(codes.InvalidArgument, "user_id is required")
		level.Error(logger).Log("msg", err.Error())
		return UpdateUserResponse{}, err
	}

	if req.HashedRefreshToken == "" || req.HashedAccessToken == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "both refresh and access tokens are required")
		level.Error(logger).Log("msg", err.Error())
		return UpdateUserResponse{}, err
	}

	_, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				level.Warn(logger).Log("msg", "user not found", "user_id", req.UserID)
				return UpdateUserResponse{}, error_hndl.NewError(codes.NotFound, "user not found: "+err.Error())
			case codes.InvalidArgument:
				level.Error(logger).Log("msg", "invalid argument", "err", err)
				return UpdateUserResponse{}, error_hndl.NewError(codes.InvalidArgument, "invalid user ID: "+err.Error())
			default:
				level.Error(logger).Log("msg", "unexpected error getting user", "err", err)
				return UpdateUserResponse{}, error_hndl.NewError(codes.Internal, "unexpected error getting user: "+err.Error())
			}
		}

		level.Error(logger).Log("msg", "unknown error getting user", "err", err)
		return UpdateUserResponse{}, error_hndl.NewError(codes.Internal, "failed to verify user existence: "+err.Error())
	}

	err = s.repo.UserRepo.UpdateUserTokens(req.UserID, req.HashedRefreshToken, req.HashedAccessToken)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				level.Warn(logger).Log("msg", "invalid token format", "err", err)
				return UpdateUserResponse{}, error_hndl.NewError(codes.InvalidArgument, "invalid token format: "+err.Error())
			case codes.Aborted:
				level.Warn(logger).Log("msg", "update aborted", "err", err)
				return UpdateUserResponse{}, error_hndl.NewError(codes.Aborted, "update aborted: "+err.Error())
			default:
				level.Error(logger).Log("msg", "failed to update tokens", "err", err)
				return UpdateUserResponse{}, error_hndl.NewError(codes.Internal, "failed to update tokens: "+err.Error())
			}
		}

		level.Error(logger).Log("msg", "unknown error updating tokens", "err", err)
		return UpdateUserResponse{}, error_hndl.NewError(codes.Internal, "unknown error updating tokens: "+err.Error())
	}

	level.Info(logger).Log("msg", "tokens updated successfully", "user_id", req.UserID)

	return UpdateUserResponse{
		Success: true,
		Message: "tokens updated successfully",
	}, nil
}
