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

	Jwt struct {
		Secret          string `mapstructure:"secret"`
		AccessTokenTtl  string `mapstructure:"access_token_ttl"`
		RefreshTokenTtl string `mapstructure:"refresh_token_ttl"`
		HashPepper      string `mapstructure:"hash_pepper"`
	} `mapstructure:"jwt"`

	AuthService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"auth_service"`

	UserDatabase struct {
		URL            string `mapstructure:"url"`
		MigrationsPath string `mapstructure:"migrations_path"`
	} `mapstructure:"user_database"`

	Tracing struct {
		Enabled      bool   `mapstructure:"enabled"`
		OtlpEndpoint string `mapstructure:"otlp_endpoint"`
	} `mapstructure:"tracing"`

	Metrics struct {
		Enabled        bool   `mapstructure:"enabled"`
		OtlpEndpoint   string `mapstructure:"otlp_endpoint"`
		ExportInterval string `mapstructure:"export_interval"`
		ExportTimeout  string `mapstructure:"export_timeout"`
	} `mapstructure:"metrics"`
}
