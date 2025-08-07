package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"customapp/internal/api"
	"customapp/internal/logger"
)

type Config struct {
	HttpServer api.HttpServer
	Logger     logger.Config
}

func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
