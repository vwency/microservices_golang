package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/hello_service/endpoint"
	"github.com/vwency/microservices_golang/internal/hello_service/service"
	"github.com/vwency/microservices_golang/internal/hello_service/transport"
	"github.com/vwency/microservices_golang/pkg/config"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "hello_service", &Cfg)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "service", Cfg.App.ServiceName)
	}

	svc := service.NewHelloService(logger)
	endpoints := endpoint.MakeEndpoints(svc)
	grpcServerImpl := transport.NewGRPCServer(endpoints, logger)

	listener, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to listen gRPC", "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	transport.RegisterHelloServiceServer(grpcServer, grpcServerImpl)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("msg", "starting gRPC server", "address", Cfg.App.Port)
		errs <- grpcServer.Serve(listener)
	}()

	level.Info(logger).Log("exit", <-errs)
}
