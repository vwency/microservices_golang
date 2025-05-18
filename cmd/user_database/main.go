package main

import (
	gokit "github.com/vwency/microservices_golang/cmd/user_database/go-kit"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			loadConfig,
			NewZapLogger,
			NewKitLogger,
			newDatabaseConnection,
			newRepository,
			gokit.NewService,
			gokit.NewEndpoints,
			gokit.NewGRPCServer,
			NewTLSCredentials,
			newListener,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}
