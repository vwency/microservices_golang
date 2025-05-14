package service

import (
	"context"
	"errors"

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

var user *models.User
var err error

func (s *userService) GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	if req.Username == "" && req.Email == "" && req.UserID == "" {
		return GetUserResponse{}, status.Error(codes.InvalidArgument, "username, email or userID must be provided")
	}

	if req.UserID != "" {
		user, err = s.repo.UserRepo.GetUserByID(req.UserID)
	} else {
		user, err = s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	}

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Log("error", "user not found", "username", req.Username, "email", req.Email, "userID", req.UserID)
			return GetUserResponse{Found: false}, nil
		}
		s.logger.Log("error", "database error", "err", err)
		return GetUserResponse{}, status.Error(codes.Internal, "internal server error")
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

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
