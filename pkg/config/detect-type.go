package config

import "os"

func DetectEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	return env
}
