package service

import (
	"context"
	stdErrors "errors" // Стандартный пакет errors с псевдонимом
	"net"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/vwency/microservices_golang/internal/auth_service/service/errors" // твой кастомный errors
	"google.golang.org/grpc/codes"
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
		err := errors.NewError(codes.InvalidArgument, "config path is required")
		level.Error(logger).Log("msg", err.Error())
		return InitDatabaseResponse{Success: false}, err
	}

	err := s.repo.UserRepo.RunMigrations()
	if err != nil {
		level.Error(logger).Log(
			"msg", "database initialization failed",
			"configPath", req.ConfigPath,
			"err", err,
		)

		if isNetworkError(err) {
			errUnavailable := errors.NewError(codes.Unavailable, "database is unavailable: "+err.Error())
			return InitDatabaseResponse{Success: false}, errUnavailable
		}

		var pathErr *os.PathError
		if stdErrors.As(err, &pathErr) { // Используем stdErrors.As для проверки типа ошибки
			errInvalid := errors.NewError(codes.InvalidArgument, "invalid config path: "+err.Error())
			return InitDatabaseResponse{Success: false}, errInvalid
		}

		errInternal := errors.NewError(codes.Internal, "failed to initialize database: "+err.Error())
		return InitDatabaseResponse{Success: false}, errInternal
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
	return stdErrors.As(err, &netErr) // Здесь тоже используем stdErrors.As
}
