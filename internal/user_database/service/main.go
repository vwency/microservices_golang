package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
)

type Service interface {
	InitDatabase(ctx context.Context, req InitDatabaseRequest) (InitDatabaseResponse, error)
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

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &userService{
		repo:   repo,
		logger: log.With(logger, "service", "user"),
	}
}
