package service

import (
	"context"

	update_user "github.com/vwency/microservices_golang/internal/user_database/service/update_user"
)

type UpdateUserRequest = update_user.Request
type UpdateUserResponse = update_user.Response

func (s *userService) UpdateUser(ctx context.Context, req update_user.Request) (update_user.Response, error) {
	return (&update_user.Service{
		Logger: s.logger,
		Repo:   &s.repo,
	}).UpdateUser(ctx, req)
}
