package main

import (
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			loadConfig,
			newZapLogger,
			newKitLogger,
			newJWTManager,
			newDatabaseConnection,
			newDatabaseClient,
			newAuthService,
			newAuthEndpoints,
			newGRPCServer,
			newListener,
		),
		fx.Invoke(startServer),
	)

	app.Run()
}
