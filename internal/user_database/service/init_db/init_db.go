package init_db

import (
	"context"
	"errors"
	"net"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	error_hndl "github.com/vwency/microservices_golang/internal/user_database/service/errors"
	"google.golang.org/grpc/codes"
)

type Request struct {
	ConfigPath string
}

type Response struct {
	Success bool
}

type Service struct {
	Logger log.Logger
	Repo   interface {
		RunMigrations() error
	}
}

func (s *Service) InitDatabase(ctx context.Context, req Request) (Response, error) {
	logger := log.With(s.Logger, "method", "InitDatabase")

	if req.ConfigPath == "" {
		err := error_hndl.NewError(codes.InvalidArgument, "config path is required")
		level.Error(logger).Log("msg", err.Error())
		return Response{Success: false}, err
	}

	err := s.Repo.RunMigrations()
	if err != nil {
		level.Error(logger).Log(
			"msg", "database initialization failed",
			"configPath", req.ConfigPath,
			"err", err,
		)

		if isNetworkError(err) {
			errUnavailable := error_hndl.NewError(codes.Unavailable, "database is unavailable: "+err.Error())
			return Response{Success: false}, errUnavailable
		}

		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			errInvalid := error_hndl.NewError(codes.InvalidArgument, "invalid config path: "+err.Error())
			return Response{Success: false}, errInvalid
		}

		errInternal := error_hndl.NewError(codes.Internal, "failed to initialize database: "+err.Error())
		return Response{Success: false}, errInternal
	}

	level.Info(logger).Log(
		"msg", "database initialized successfully",
		"configPath", req.ConfigPath,
	)

	return Response{Success: true}, nil
}

func isNetworkError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}
