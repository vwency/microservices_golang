package main

import (
	"github.com/vwency/microservices_golang/pkg/config"
)

func loadConfig() config.ServiceConfig {
	env := config.DetectEnv()
	var cfg config.ServiceConfig
	config.Init(env, "user_database", &cfg)
	return cfg
}
