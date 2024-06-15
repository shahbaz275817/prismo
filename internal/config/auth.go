package config

import cfg "github.com/shahbaz275817/prismo/pkg/config"

type AuthConfig struct {
	Username string
	Password string
}

func newAuthConfig() AuthConfig {
	return AuthConfig{
		Username: cfg.MustGetString("AUTH_USERNAME"),
		Password: cfg.MustGetString("AUTH_PASSWORD"),
	}
}
