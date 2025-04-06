package config

type ServiceConfig struct {
	App struct {
		Env         string `mapstructure:"env"`
		Port        string `mapstructure:"port"`
		LogLevel    string `mapstructure:"log_level"`
		ServiceName string `mapstructure:"service_name"`
	} `mapstructure:"app"`

	Database struct {
		URL            string `mapstructure:"url"`
		MigrationsPath string `mapstructure:"migrations_path"`
	} `mapstructure:"database"`
}
