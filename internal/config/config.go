package config

import "time"

type (
	Config struct {
		Server        ServerConfig        `yaml:"server"`
		Redis         RedisConfig         `yaml:"redis"`
		CurrencyCache CurrencyCacheConfig `yaml:"currency_cache"`
		Provider      ProviderConfig      `yaml:"provider"`
	}

	ServerConfig struct {
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	}

	RedisConfig struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		Password     string        `env:"REDIS_PASSWORD" env-required:"true"`
		MinIdleConns int           `yaml:"min_idle_conns"`
		MaxIdleConns int           `yaml:"max_idle_conns"`
		DialTimeout  time.Duration `yaml:"dial_timeout"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	}

	CurrencyCacheConfig struct {
		BaseTTL      time.Duration `yaml:"base_ttl"`
		StaleTTL     time.Duration `yaml:"stale_ttl"`
		RefreshDelta time.Duration `yaml:"refresh_delta"`
	}

	ProviderConfig struct {
		EasyBit EasyBitConfig `yaml:"easybit"`
	}

	EasyBitConfig struct {
		BaseURL             string        `yaml:"base_url"`
		APIKey              string        `env:"EASYBIT_API_KEY" env-required:"true"`
		LimiterRate         int           `yaml:"limiter_rate"`
		LimiterBurst        int           `yaml:"limiter_burst"`
		AttemptTimeout      time.Duration `yaml:"attempt_timeout"`
		MaxRetries          uint          `yaml:"max_retries"`
		RequestTimeout      time.Duration `yaml:"request_timeout"`
		MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	}
)
