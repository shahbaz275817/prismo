package config

import (
	"fmt"

	"github.com/shahbaz275817/prismo/pkg/cache"
	cfg "github.com/shahbaz275817/prismo/pkg/config"
	"github.com/shahbaz275817/prismo/pkg/locks"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

var appConfig config

type config struct {
	app              AppConfig
	logger           logger.Config
	db               DBConfig
	readDB           DBConfig
	cache            cache.Options
	auth             AuthConfig
	atomicLockConfig map[locks.KeyType]locks.LockConfig
}

func Load() {
	cfg.Init()
	initDBConfig()
	initStatsDConfig()

	appConfig = config{
		app:              newAppConfig(),
		logger:           logger.NewConfig(),
		db:               dbConf(),
		readDB:           readDBConf(),
		cache:            cache.NewCacheConfig(),
		auth:             newAuthConfig(),
		atomicLockConfig: newAtomicLockConfig(),
	}
}

func Port() string                                         { return fmt.Sprintf("%d", appConfig.app.Port) }
func Addr() string                                         { return fmt.Sprintf("%s:%d", appConfig.app.Host, appConfig.app.Port) }
func Log() logger.Config                                   { return appConfig.logger }
func Auth() AuthConfig                                     { return appConfig.auth }
func DB() DBConfig                                         { return appConfig.db }
func ReadDB() DBConfig                                     { return appConfig.readDB }
func Cache() cache.Options                                 { return appConfig.cache }
func AtomicLockConfig() map[locks.KeyType]locks.LockConfig { return appConfig.atomicLockConfig }
