package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserRequest struct {
	UserID   string
	Username string
	Email    string
}

type GetUserResponse struct {
	Found              bool
	UserID             string
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
}

func (s *userService) GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if req.Username == "" && req.Email == "" && req.UserID == "" {
		level.Error(logger).Log("msg", "no search criteria provided")
		return GetUserResponse{}, NewInvalidArgumentError("username, email or userID must be provided", nil)
	}

	var user *models.User
	var err error

	if req.UserID != "" {
		user, err = s.repo.UserRepo.GetUserByID(req.UserID)
	} else {
		user, err = s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	}

	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			level.Debug(logger).Log(
				"msg", "user not found",
				"username", req.Username,
				"email", req.Email,
				"userID", req.UserID,
			)
			// В случае, если пользователь не найден, генерируем ошибку
			return GetUserResponse{}, NewNotFoundError("user not found", err)

		case status.Code(err) == codes.InvalidArgument:
			level.Warn(logger).Log(
				"msg", "invalid request",
				"err", err,
			)
			return GetUserResponse{}, NewInvalidArgumentError("invalid request parameters", err)

		default:
			level.Error(logger).Log(
				"msg", "database error",
				"err", err,
			)
			return GetUserResponse{}, NewInternalError("failed to get user", err)
		}
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	level.Info(logger).Log(
		"msg", "user found",
		"userID", user.ID.String(),
	)

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
