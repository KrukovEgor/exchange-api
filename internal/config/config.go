package config

import "time"

type (
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Redis    RedisConfig    `yaml:"redis"`
		Provider ProviderConfig `yaml:"provider"`
	}

	ServerConfig struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	}

	RedisConfig struct {
		Port         string `yaml:"port"`
		Password     string `env:"REDIS_PASSWORD" env-required:"true"`
		MinIdleConns int    `yaml:"min_idle_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	}

	ProviderConfig struct {
		EasyBit EasyBitConfig `yaml:"easybit"`
	}

	EasyBitConfig struct {
		BaseURL             string        `yaml:"base_url"`
		APIKey              string        `env:"EASYBIT_API_KEY" env-required:"true"`
		LimiterRate         int           `yaml:"limiter_rate"`
		LimiterBurst        int           `yaml:"limiter_burst"`
		RequestTimeout      time.Duration `yaml:"request_timeout"`
		MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	}
)
