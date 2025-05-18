package main

import (
	"context"
	"fmt"
	"net"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/vwency/microservices_golang/pkg/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func StartServer(
	server *grpc.Server,
	listener net.Listener,
	logger log.Logger,
	cfg config.ServiceConfig,
	lc fx.Lifecycle,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			level.Info(logger).Log(
				"msg", "Запуск gRPC сервера auth_service",
				"env", config.DetectEnv(),
				"addr", fmt.Sprintf(":%s", cfg.App.Port),
			)

			go func() {
				if err := server.Serve(listener); err != nil {
					level.Error(logger).Log("msg", "Ошибка при запуске gRPC сервера", "err", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			level.Info(logger).Log("msg", "Остановка gRPC сервера auth_service")
			server.GracefulStop()
			return nil
		},
	})
}
