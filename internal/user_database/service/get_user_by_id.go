package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

type GetUserByIDRequest struct {
	UserID string
}

type GetUserByIDResponse struct {
	Found              bool
	UserID             string
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
}

func (s *userService) GetUserByID(ctx context.Context, req GetUserByIDRequest) (GetUserByIDResponse, error) {
	logger := log.With(s.logger, "method", "GetUserByID")

	if req.UserID == "" {
		level.Error(logger).Log("msg", "userID is required")
		return GetUserByIDResponse{}, NewInvalidArgumentError("userID must be provided", nil)
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		level.Warn(logger).Log("msg", "invalid userID format", "userID", req.UserID, "err", err)
		return GetUserByIDResponse{}, NewInvalidArgumentError("invalid userID format", err)
	}

	user, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		var serviceErr *ServiceError
		if errors.As(err, &serviceErr) {
			switch serviceErr.Code {
			case "not_found":
				level.Debug(logger).Log("msg", "user not found", "userID", req.UserID)
				return GetUserByIDResponse{Found: false}, nil

			case "invalid_argument":
				level.Warn(logger).Log("msg", "invalid userID", "userID", req.UserID, "err", err)
				return GetUserByIDResponse{}, NewInvalidArgumentError("invalid userID", err)

			default:
				level.Error(logger).Log("msg", "failed to get user", "userID", req.UserID, "err", err)
				return GetUserByIDResponse{}, NewInternalError("failed to get user", err)
			}
		}

		level.Error(logger).Log("msg", "unexpected error when getting user", "userID", req.UserID, "err", err)
		return GetUserByIDResponse{}, NewInternalError("unexpected error when getting user", err)
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

	return GetUserByIDResponse{
		Found:              true,
		UserID:             userID,
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
	}, nil
}
