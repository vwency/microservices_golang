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

// GetUser retrieves a user by username or email.

// GetUserByID retrieves a user by their user ID.
func (s *userService) GetUserByID(ctx context.Context, req GetUserByIDRequest) (GetUserByIDResponse, error) {
	// Validate request
	if req.UserID == "" {
		return GetUserByIDResponse{}, status.Error(codes.InvalidArgument, "userID must be provided")
	}

	// Fetch user by ID (passing the string directly)
	user, err := s.repo.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		// Log the error and return a gRPC NotFound error
		s.logger.Log("error", "user not found", "userID", req.UserID)
		return GetUserByIDResponse{}, status.Error(codes.NotFound, "user not found")
	}

	// Handle pointer fields and type conversions
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	// Convert UUID to string if needed
	userID := req.UserID
	if user.ID != uuid.Nil { // If the repository returned a UUID
		userID = user.ID.String()
	}

	// Return the found user
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
