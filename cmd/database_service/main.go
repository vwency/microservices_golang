package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	"github.com/vwency/microservices_golang/internal/user_database/repository/user_repository"
	"github.com/vwency/microservices_golang/internal/user_database/service"
	grpctransport "github.com/vwency/microservices_golang/internal/user_database/transport"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/database"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "user_database", &Cfg)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	db, err := database.NewGORM(Cfg.UserDatabase.URL)
	if err != nil {
		level.Error(logger).Log("msg", "failed to connect to user_database", "err", err)
		os.Exit(1)
	}

	if err := user_repository.RunUserMigrations(db); err != nil {
		level.Error(logger).Log("msg", "failed to run migrations", "err", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(db)

	var svc service.Service
	{
		svc = service.NewService(*repo, logger)
	}

	eps := endpoints.MakeEndpoints(svc)

	grpcListener, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to listen", "err", err)
		os.Exit(1)
	}

	baseServer := grpc.NewServer(
		grpc.UnaryInterceptor(kitgrpc.Interceptor),
	)

	grpctransport.RegisterGRPCServer(baseServer, eps)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log(
			"msg", "starting gRPC server",
			"service", Cfg.App.ServiceName,
			"port", Cfg.App.Port,
		)
		errs <- baseServer.Serve(grpcListener)
	}()

	level.Info(logger).Log("msg", "exiting", "err", <-errs)
}
