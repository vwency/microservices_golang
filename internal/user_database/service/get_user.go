package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
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
	Found              bool   `json:"found"`
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	HashedPassword     string `json:"hashed_password"`
	HashedRefreshToken string `json:"hashed_refresh_token"`
	HashedAccessToken  string `json:"hashed_access_token"`
}

func (s *userService) GetUser(ctx context.Context, request GetUserRequest) (GetUserResponse, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if request.UserID == "" && request.Username == "" && request.Email == "" {
		level.Error(logger).Log("msg", "no search criteria provided")
		return GetUserResponse{}, NewInvalidArgumentError("username, email or userID must be provided", nil)
	}

	var (
		user *models.User
		err  error
	)

	if request.UserID != "" {
		_, parseErr := uuid.Parse(request.UserID)
		if parseErr != nil {
			level.Warn(logger).Log("msg", "invalid user_id format", "user_id", request.UserID, "err", parseErr)
			return GetUserResponse{}, NewInvalidArgumentError("invalid user_id format", parseErr)
		}

		user, err = s.repo.UserRepo.GetUserByID(request.UserID)
	} else {
		user, err = s.repo.UserRepo.GetUserByUsernameOrEmail(request.Username, request.Email)
	}

	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			level.Warn(logger).Log("msg", "user not found", "user_id", request.UserID, "username", request.Username, "email", request.Email)
			return GetUserResponse{}, NewNotFoundError("user not found", err)
		case codes.InvalidArgument:
			level.Warn(logger).Log("msg", "invalid request parameters", "err", err)
			return GetUserResponse{}, NewInvalidArgumentError("invalid request parameters", err)
		default:
			level.Error(logger).Log("msg", "failed to get user", "err", err)
			return GetUserResponse{}, NewInternalError("failed to get user", err)
		}
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	level.Info(logger).Log("msg", "user found", "user_id", user.ID.String(), "username", user.Username)

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
