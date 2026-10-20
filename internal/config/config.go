package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port      int    `env:"PORT"       envDefault:"8080"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("config: PORT out of range: %d", c.Port)
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("config: JWT_SECRET must be at least 32 chars")
	}

	return nil
}
