package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
)

type GetUserRequest struct {
	UserID   string
	Username string
	Email    string
}

type GetUserResponse struct {
	Found              bool   `json:"found"`
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	HashedPassword     string `json:"hashed_password"`
	HashedRefreshToken string `json:"hashed_refresh_token"`
	HashedAccessToken  string `json:"hashed_access_token"`
}

func (s *userService) GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if req.UserID == "" && req.Username == "" && req.Email == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "username, email or userID must be provided")
		level.Error(logger).Log("err", err)
		return GetUserResponse{}, err
	}

	var user *models.User
	var err error

	if req.UserID != "" {
		if _, parseErr := uuid.Parse(req.UserID); parseErr != nil {
			err = error_hndl.NewError(codes.InvalidArgument, "invalid user_id format")
			level.Warn(logger).Log("err", err, "user_id", req.UserID)
			return GetUserResponse{}, err
		}

		user, err = s.repo.UserRepo.GetUserByID(req.UserID)
	} else {
		user, err = s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	}

	if err != nil {
		var grpcErr *error_hndl.Error
		if error_hndl.As(err, &grpcErr) {
			level.Warn(logger).Log("err", grpcErr)
			return GetUserResponse{}, grpcErr
		}

		if error_hndl.Is(err, error_hndl.ErrNotFound) {
			err = error_hndl.NewError(codes.NotFound, "user not found")
			level.Warn(logger).Log("err", err)
			return GetUserResponse{}, err
		}

		err = error_hndl.NewError(codes.Internal, "failed to get user")
		level.Error(logger).Log("err", err)
		return GetUserResponse{}, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	level.Info(logger).Log("msg", "user found", "userID", user.ID.String())

	return GetUserResponse{
		Found:              true,
		UserID:             user.ID.String(),
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
	}, nil
}
