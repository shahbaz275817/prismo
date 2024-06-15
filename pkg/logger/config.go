package logger

import "github.com/shahbaz275817/prismo/pkg/config"

type Config struct {
	LogLevel string
	Format   string
}

func NewConfig() Config {
	return Config{
		LogLevel: config.MustGetString("LOG_LEVEL"),
		Format:   config.GetString("LOG_FORMAT"),
	}
}
