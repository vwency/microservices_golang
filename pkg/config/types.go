package config

// ServiceConfig содержит настройки для сервиса
type ServiceConfig struct {
	App struct {
		Env      string `mapstructure:"env"`
		Port     string `mapstructure:"port"`
		LogLevel string `mapstructure:"log_level"`
	}
}
