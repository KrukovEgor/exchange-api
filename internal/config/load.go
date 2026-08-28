package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func Load(configPath string) (*Config, error) {
	const op = "config.Load"

	var cfg Config

	if configPath == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to read env variables: %w", op, err)
		}

		return &cfg, nil
	}

	_ = godotenv.Load()

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %w", op, err)
	}

	return &cfg, nil
}
