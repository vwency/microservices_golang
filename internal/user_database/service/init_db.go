package service

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type InitDatabaseRequest struct {
	ConfigPath string
}

type InitDatabaseResponse struct {
	Success bool
}

func (s *userService) InitDatabase(ctx context.Context, req InitDatabaseRequest) (InitDatabaseResponse, error) {
	logger := log.With(s.logger, "method", "InitDatabase")

	// Validate request
	if req.ConfigPath == "" {
		level.Error(logger).Log("msg", "config path is required")
		return InitDatabaseResponse{Success: false}, NewInvalidArgumentError("config path is required", nil)
	}

	// Run database migrations
	err := s.repo.UserRepo.RunMigrations()
	if err != nil {
		level.Error(logger).Log(
			"msg", "database initialization failed",
			"configPath", req.ConfigPath,
			"err", err,
		)

		if isConnectionError(err) {
			return InitDatabaseResponse{Success: false}, NewUnavailableError("database unavailable", err)
		}
		return InitDatabaseResponse{Success: false}, NewInternalError("failed to initialize database", err)
	}

	level.Info(logger).Log(
		"msg", "database initialized successfully",
		"configPath", req.ConfigPath,
	)

	return InitDatabaseResponse{
		Success: true,
	}, nil
}

func isConnectionError(err error) bool {
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "host unreachable")
}
