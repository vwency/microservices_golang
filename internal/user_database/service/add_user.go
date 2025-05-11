package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	RunMigrations() error
	GetUserByID(id string) (*models.User, error)
	GetUserByUsernameOrEmail(username, email string) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUserTokens(userID, hashedRefreshToken, hashedAccessToken string) error
}

// Внутри пакета service
type AddUserRequest struct {
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
	UserID             string
}

type AddUserResponse struct {
	Success bool
	Message string
}

func strPtr(s string) *string {
	return &s
}

func (s *UserService) AddUserRequest(ctx context.Context, req AddUserRequest) (*AddUserResponse, error) {
	user := &models.User{
		Username:           req.Username,
		Email:              strPtr(req.Email),
		HashedPassword:     req.HashedPassword,
		HashedRefreshToken: req.HashedRefreshToken,
		HashedAccessToken:  req.HashedAccessToken,
	}

	if req.UserID != "" {
		user.ID, _ = uuid.Parse(req.UserID) // если ID задан, парсим UUID
	}

	// Check if user exists
	existingUser, err := s.repo.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return &AddUserResponse{
			Success: false,
			Message: "user already exists",
		}, nil
	}

	// Add new user
	err = s.repo.AddUser(user)
	if err != nil {
		return nil, err
	}

	return &AddUserResponse{
		Success: true,
		Message: "user created successfully",
	}, nil
}
