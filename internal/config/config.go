package config

type (
	Config struct {
		Server ServerConfig `yaml:"server"`
	}

	ServerConfig struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	}
)
