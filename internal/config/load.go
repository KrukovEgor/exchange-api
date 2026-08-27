package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load(configPath string) (*Config, error) {
	const op = "Load"

	var cfg Config

	if configPath == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to read env variables: %w", op, err)
		}

		return &cfg, nil
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %w", op, err)
	}

	return &cfg, nil
}
