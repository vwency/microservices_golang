package get_user

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
)

type Service struct {
	Logger log.Logger
	Repo   *repository.Repository
}

type Request struct {
	UserID   string
	Username string
	Email    string
}

type Response struct {
	Found              bool   `json:"found"`
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	HashedPassword     string `json:"hashed_password"`
	HashedRefreshToken string `json:"hashed_refresh_token"`
	HashedAccessToken  string `json:"hashed_access_token"`
}

func (s *Service) GetUser(ctx context.Context, req Request) (Response, error) {
	logger := log.With(s.Logger, "method", "GetUser")

	if req.UserID == "" && req.Username == "" && req.Email == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "username, email or userID must be provided")
		level.Error(logger).Log("err", err)
		return Response{}, err
	}

	var user *models.User
	var err error

	if req.UserID != "" {
		if _, parseErr := uuid.Parse(req.UserID); parseErr != nil {
			err = error_hndl.NewError(codes.InvalidArgument, "invalid user_id format")
			level.Warn(logger).Log("err", err, "user_id", req.UserID)
			return Response{}, err
		}

		user, err = s.Repo.UserRepo.GetUserByID(req.UserID)
	} else {
		user, err = s.Repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	}

	if err != nil {
		var grpcErr *error_hndl.Error
		if error_hndl.As(err, &grpcErr) {
			level.Warn(logger).Log("err", grpcErr)
			return Response{}, grpcErr
		}

		if error_hndl.Is(err, error_hndl.ErrNotFound) {
			err = error_hndl.NewError(codes.NotFound, "user not found")
			level.Warn(logger).Log("err", err)
			return Response{}, err
		}

		err = error_hndl.NewError(codes.Internal, "failed to get user")
		level.Error(logger).Log("err", err)
		return Response{}, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	level.Info(logger).Log("msg", "user found", "userID", user.ID.String())

	return Response{
		Found:              true,
		UserID:             user.ID.String(),
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
	}, nil
}
