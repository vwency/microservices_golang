package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
)

type Service interface {
	AddUser(ctx context.Context, request AddUserRequest) (AddUserResponse, error)
	GetUser(ctx context.Context, request GetUserRequest) (GetUserResponse, error)
	GetUserByID(ctx context.Context, request GetUserByIDRequest) (GetUserByIDResponse, error)
	UpdateUser(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error)
	DeleteUser(ctx context.Context, request DeleteUserRequest) (DeleteUserResponse, error)
}

type userService struct {
	repo   repository.Repository
	logger log.Logger
}

// AddUser implements Service.
func (s *userService) AddUser(ctx context.Context, request AddUserRequest) (AddUserResponse, error) {
	panic("unimplemented")
}

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &userService{
		repo:   repo,
		logger: log.With(logger, "service", "user"),
	}
}
