package get_user_by_id

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vwency/microservices_golang/internal/user_database/models"
)

type Request struct {
	UserID string
}

type Response struct {
	Found              bool
	UserID             string
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
}

// интерфейс репозитория должен совпадать с тем, что реализует user_repository.UserRepository
type UserRepo interface {
	GetUserByID(id string) (*models.User, error)
}

type Service struct {
	Logger log.Logger
	Repo   UserRepo
}

func (s *Service) GetUserByID(ctx context.Context, req Request) (Response, error) {
	logger := log.With(s.Logger, "method", "GetUserByID")

	if req.UserID == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "userID must be provided")
		level.Error(logger).Log("msg", err.Error())
		return Response{}, err
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		errInvalid := error_hndl.NewError(codes.InvalidArgument, "invalid userID format: "+err.Error())
		level.Warn(logger).Log("msg", errInvalid.Error(), "userID", req.UserID, "err", err)
		return Response{}, errInvalid
	}

	user, err := s.Repo.GetUserByID(req.UserID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			level.Debug(logger).Log("msg", "user not found", "userID", req.UserID)
			return Response{Found: false}, nil
		}

		errInternal := error_hndl.NewError(codes.Internal, "failed to get user: "+err.Error())
		level.Error(logger).Log("msg", errInternal.Error(), "userID", req.UserID, "err", err)
		return Response{}, errInternal
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	userID := req.UserID
	if user.ID != uuid.Nil {
		userID = user.ID.String()
	}

	level.Debug(logger).Log("msg", "user found", "userID", userID)

	return Response{
		Found:              true,
		UserID:             userID,
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
	}, nil
}
