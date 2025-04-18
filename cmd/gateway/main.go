package main

import (
	"log"
	"net/http"

	"github.com/vwency/microservices_golang/internal/gateway"
	"github.com/vwency/microservices_golang/internal/gateway/service"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/logger"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "gateway", &Cfg)

	logger.Init(Cfg.App.LogLevel)

	authService, err := service.NewAuthServiceClient(Cfg.AuthService.URL)

	if err != nil {
		logger.Fatal("could not create service clients: %v", err)
	}

	r := gateway.InitializeRouter(authService)

	logger.Info("Starting HTTP server on port " + Cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+Cfg.App.Port, r))
}
