package config

import (
	cfg "github.com/shahbaz275817/prismo/pkg/config"
)

type AppConfig struct {
	Host string
	Port int
}

func newAppConfig() AppConfig {
	return AppConfig{
		Host: cfg.MustGetString("APP_HOST"),
		Port: cfg.MustGetInt("APP_PORT"),
	}
}
