package config

import (
	"os"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
	LogLevel   string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBUrl:      os.Getenv("DB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	return cfg
}
