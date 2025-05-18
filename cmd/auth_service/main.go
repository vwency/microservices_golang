package main

import (
	gokit "github.com/vwency/microservices_golang/cmd/auth_service/go-kit"
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
			gokit.NewAuthService,
			gokit.NewAuthEndpoints,
			gokit.NewGRPCServer,
			newListener,
			NewTLSCredentials,
		),
		fx.Invoke(startServer),
	)

	app.Run()
}
