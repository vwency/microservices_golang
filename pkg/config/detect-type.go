package config

import "os"

// DetectEnv определяет тип окружения (например, dev, prod)
func DetectEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Значение по умолчанию
	}
	return env
}
