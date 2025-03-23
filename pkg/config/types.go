package config

type ServiceConfig struct {
	App struct {
		Env      string `mapstructure:"env"`
		Port     string `mapstructure:"port"`
		LogLevel string `mapstructure:"log_level"`
		service_name string `mapstructure:"service_name"`
	}
}
