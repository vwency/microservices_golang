package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if req.UserID == "" {
		return GetUserByIDResponse{}, status.Error(codes.InvalidArgument, "userID must be provided")
	}

	user, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		s.logger.Log("error", "user not found", "userID", req.UserID)
		return GetUserByIDResponse{}, status.Error(codes.NotFound, "user not found")
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	userID := req.UserID
	if user.ID != uuid.Nil {
		userID = user.ID.String()
	}

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
