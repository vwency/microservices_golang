package service

import (
	"context"

	"github.com/google/uuid"
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
func (s *userService) GetUser(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	// Validate request
	if req.Username == "" && req.Email == "" && req.UserID == "" {
		return GetUserResponse{}, status.Error(codes.InvalidArgument, "username, email or userID must be provided")
	}

	// Fetch user by username or email
	user, err := s.repo.UserRepo.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		// Log the error and return a gRPC NotFound error
		s.logger.Log("error", "user not found", "username", req.Username, "email", req.Email)
		return GetUserResponse{}, status.Error(codes.NotFound, "user not found")
	}

	// Handle pointer fields and type conversions
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	// Return the found user
	return GetUserResponse{
		Found:              true,
		UserID:             user.ID.String(), // Convert UUID to string
		Username:           user.Username,
		Email:              email, // Handle nil pointer
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
	}, nil
}

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
