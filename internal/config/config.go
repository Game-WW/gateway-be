package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port           int    `env:"PORT"            envDefault:"8080"`
	JWTSecret      string `env:"JWT_SECRET,required"`
	UsersGRPCAddr  string `env:"USERS_GRPC_ADDR" envDefault:"127.0.0.1:50051"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS,required"`
	CookieSecure   bool   `env:"COOKIE_SECURE"   envDefault:"true"`
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
	if c.UsersGRPCAddr == "" {
		return fmt.Errorf("config: USERS_GRPC_ADDR is required")
	}
	if c.AllowedOrigins == "" {
		return fmt.Errorf("config: ALLOWED_ORIGINS is required")
	}
	for _, origin := range strings.Split(c.AllowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Host == "" || strings.Contains(parsed.Host, "*") ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("config: invalid ALLOWED_ORIGINS entry %q; expected an http(s) origin without a path", origin)
		}
	}

	return nil
}
