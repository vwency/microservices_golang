package main

import (
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
			NewTLSCredentials,
			newListener,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}
