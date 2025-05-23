package update_user

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
)

type Service struct {
	Logger log.Logger
	Repo   *repository.Repository
}

type Request struct {
	UserID             string
	HashedRefreshToken string
	HashedAccessToken  string
}

type Response struct {
	Success bool
	Message string
}

func (s *Service) UpdateUser(ctx context.Context, req Request) (Response, error) {
	logger := log.With(s.Logger, "method", "UpdateUser")

	if req.UserID == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "user_id is required")
		level.Error(logger).Log("msg", err.Error())
		return Response{}, err
	}

	if req.HashedRefreshToken == "" || req.HashedAccessToken == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "both refresh and access tokens are required")
		level.Error(logger).Log("msg", err.Error())
		return Response{}, err
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		err := error_hndl.NewError(codes.InvalidArgument, "invalid user_id format")
		level.Error(logger).Log("msg", err.Error())
		return Response{}, err
	}

	user, err := s.Repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if error_hndl.Is(err, error_hndl.ErrNotFound) {
			level.Warn(logger).Log("msg", "user not found", "user_id", req.UserID)
			return Response{}, error_hndl.NewError(codes.NotFound, "user not found")
		}

		level.Error(logger).Log("msg", "failed to verify user existence", "err", err)
		return Response{}, error_hndl.NewError(codes.Internal, "failed to verify user existence")
	}

	err = s.Repo.UserRepo.UpdateUserTokens(user.ID.String(), req.HashedRefreshToken, req.HashedAccessToken)
	if err != nil {
		level.Error(logger).Log("msg", "failed to update tokens", "err", err)
		return Response{}, error_hndl.NewError(codes.Internal, "failed to update tokens")
	}

	level.Info(logger).Log("msg", "tokens updated successfully", "user_id", req.UserID)
	return Response{
		Success: true,
		Message: "tokens updated successfully",
	}, nil
}
