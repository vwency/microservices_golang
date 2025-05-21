package service

import (
	"context"
	"errors"
	"net"
	"os"

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

	if req.ConfigPath == "" {
		level.Error(logger).Log("msg", "config path is required")
		return InitDatabaseResponse{Success: false}, NewInvalidArgumentError("config path is required", nil)
	}

	err := s.repo.UserRepo.RunMigrations()
	if err != nil {
		level.Error(logger).Log(
			"msg", "database initialization failed",
			"configPath", req.ConfigPath,
			"err", err,
		)

		if isNetworkError(err) {
			return InitDatabaseResponse{Success: false}, NewUnavailableError("database is unavailable", err)
		}

		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return InitDatabaseResponse{Success: false}, NewInvalidArgumentError("invalid config path", err)
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

func isNetworkError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}
