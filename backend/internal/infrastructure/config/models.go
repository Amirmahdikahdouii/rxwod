package config


type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}
