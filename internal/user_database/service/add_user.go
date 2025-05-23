package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service/add_user"
)

type AddUserRequest = add_user.Request
type AddUserResponse = add_user.Response

func (s *userService) AddUser(ctx context.Context, req AddUserRequest) (AddUserResponse, error) {
	return add_user.Execute(ctx, s.logger, s.repo.UserRepo, req)
}
