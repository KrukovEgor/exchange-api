package config

import "time"

type (
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Provider ProviderConfig `yaml:"provider"`
	}

	ServerConfig struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	}

	ProviderConfig struct {
		EasyBit EasyBitConfig `yaml:"easybit"`
	}

	EasyBitConfig struct {
		BaseURL           string        `yaml:"base_url"`
		APIKey            string        `env:"EASYBIT_API_KEY" env-required:"true"`
		RequestTimeout    time.Duration `yaml:"request_timeout"`
		LimiterRate       int           `yaml:"limiter_rate"`
		LimiterBurst      int           `yaml:"limiter_burst"`
		DisableKeepAlives bool          `yaml:"disable_keep_alives"`
	}
)
