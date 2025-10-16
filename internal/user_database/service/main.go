package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	"github.com/vwency/microservices_golang/internal/user_database/service/add_user"
	"github.com/vwency/microservices_golang/internal/user_database/service/delete_user"
)

type Service interface {
	InitDatabase(ctx context.Context, req InitDatabaseRequest) (InitDatabaseResponse, error)
	AddUser(ctx context.Context, request add_user.Request) (add_user.Response, error)
	GetUser(ctx context.Context, request GetUserRequest) (GetUserResponse, error)
	GetUserByID(ctx context.Context, request GetUserByIDRequest) (GetUserByIDResponse, error)
	UpdateUser(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error)
	DeleteUser(ctx context.Context, request delete_user.Request) (delete_user.Response, error)
}

type userService struct {
	repo          repository.Repository
	logger        log.Logger
	addUserSvc    add_user.Service
	deleteUserSvc delete_user.Service
}

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &userService{
		repo:          repo,
		logger:        log.With(logger, "service", "user"),
		addUserSvc:    add_user.NewService(repo.UserRepo, logger),
		deleteUserSvc: delete_user.NewService(repo.UserRepo, logger),
	}
}

func (s *userService) AddUser(ctx context.Context, req add_user.Request) (add_user.Response, error) {
	return s.addUserSvc.AddUser(ctx, req)
}

func (s *userService) DeleteUser(ctx context.Context, req delete_user.Request) (delete_user.Response, error) {
	return s.deleteUserSvc.DeleteUser(ctx, req)
}
