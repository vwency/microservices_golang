package service

import (
	"context"

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
	// Validate request
	if req.Username == "" && req.Email == "" && req.UserID == "" {
		return GetUserResponse{}, status.Error(codes.InvalidArgument, "username, email or userID must be provided")
	}

	user, err := s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		s.logger.Log("error", "user not found", "username", req.Username, "email", req.Email)
		return GetUserResponse{}, status.Error(codes.NotFound, "user not found")
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
