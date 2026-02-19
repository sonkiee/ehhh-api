package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppEnv  string
	Port    string
	DBURL   string
	Timeout string
}

func LoadConfig() Config {
	timeout := 5
	if v := os.Getenv("REQUEST_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeout = n
		}
	}

	cfg := Config{
		AppEnv:  os.Getenv("APP_ENV"),
		Port:    os.Getenv("PORT"),
		DBURL:   os.Getenv("DATABASE_URL"),
		Timeout: strconv.Itoa(timeout) + "s",
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	return cfg
}
