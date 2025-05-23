package service

import (
	"context"

	get_user_by_id "github.com/vwency/microservices_golang/internal/user_database/service/get_user_by_id"
)

type GetUserByIDRequest = get_user_by_id.Request
type GetUserByIDResponse = get_user_by_id.Response

func (s *userService) GetUserByID(ctx context.Context, req GetUserByIDRequest) (GetUserByIDResponse, error) {
	return (&get_user_by_id.Service{
		Logger: s.logger,
		Repo:   s.repo.UserRepo,
	}).GetUserByID(ctx, req)
}
