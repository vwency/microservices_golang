package service

import (
	"context"

	init_db "github.com/vwency/microservices_golang/internal/user_database/service/init_db"
)

type InitDatabaseRequest = init_db.Request
type InitDatabaseResponse = init_db.Response

func (s *userService) InitDatabase(ctx context.Context, req InitDatabaseRequest) (InitDatabaseResponse, error) {
	return (&init_db.Service{
		Logger: s.logger,
		Repo:   s.repo.UserRepo,
	}).InitDatabase(ctx, req)
}
