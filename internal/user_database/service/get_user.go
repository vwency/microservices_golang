package service

import (
	"context"

	get_user "github.com/vwency/microservices_golang/internal/user_database/service/get_user"
)

type GetUserRequest = get_user.Request
type GetUserResponse = get_user.Response

func (s *userService) GetUser(ctx context.Context, req get_user.Request) (get_user.Response, error) {
	return (&get_user.Service{
		Logger: s.logger,
		Repo:   &s.repo,
	}).GetUser(ctx, req)
}
