package server

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BindAddr      string      `toml:"bind_addr"`
	LogLevel      string      `toml:"log_level"`
	LogFormat     string      `toml:"log_format"`
	RedisAddr     string      `toml:"redis_addr"`
	RedisUsername string      `toml:"redis_username"`
	RedisPassword string      `toml:"redis_password"`
	Store         StoreConfig `toml:"store"`
}

type StoreConfig struct {
	DatabaseURL string `toml:"database_url"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
