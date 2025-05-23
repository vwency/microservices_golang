package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/service/delete_user"
)

type DeleteUserRequest = delete_user.Request
type DeleteUserResponse = delete_user.Response

func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	return delete_user.Execute(ctx, s.logger, s.repo.UserRepo, req)
}
