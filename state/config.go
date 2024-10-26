package state

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	ApplicationPort int     `env:"APPLICATION_PORT" envDefault:""`
	DatabaseUrl     string  `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/contacts?sslmode=disable"`
	LogLevel        string  `env:"LOG_LEVEL" envDefault:"debug"`
	SecretKey       string  `env:"SECRET_KEY" envDefault:"my_jwt_secret"`
	Rps             float64 `env:"limiter_rps" envDefault:"0"`
	Burst           int     `env:"limiter_burst" envDefault:"0"`
	LimiterEnabled  bool    `env:"limiter_enabled" envDefault:"false"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.ParseWithOptions(cfg, env.Options{RequiredIfNoDef: true})
	return cfg, err
}
