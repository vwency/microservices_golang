package service

import (
	"context"

	"github.com/go-kit/kit/log/level"
)

type InitDatabaseRequest struct {
	ConfigPath string
}

type InitDatabaseResponse struct {
	Success bool
}

func (s *userService) InitDatabase(ctx context.Context, req InitDatabaseRequest) (InitDatabaseResponse, error) {
	if req.ConfigPath == "" {
		return InitDatabaseResponse{Success: false}, ErrInvalidArgument
	}

	err := s.repo.UserRepo.RunMigrations()
	if err != nil {
		level.Error(s.logger).Log("msg", "database initialization failed", "configPath", req.ConfigPath, "err", err)
		return InitDatabaseResponse{Success: false}, ErrInternal
	}

	level.Info(s.logger).Log("msg", "database initialized successfully", "configPath", req.ConfigPath)
	return InitDatabaseResponse{Success: true}, nil
}
