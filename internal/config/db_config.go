package config

import (
	"fmt"

	configUtil "github.com/shahbaz275817/prismo/pkg/config"
)

var dbCfg DBConfig
var readDBCfg DBConfig

func dbConf() DBConfig {
	return dbCfg
}

func readDBConf() DBConfig {
	return readDBCfg
}

type DBConfig struct {
	host                   string
	port                   string
	name                   string
	user                   string
	password               string
	sslmode                string
	maxPoolSize            int
	maxIdleConnections     int
	connMaxLifetimeMinutes int
	transactionTimeoutInMS int
	isTest                 bool
}

func initDBConfig() {
	dbCfg = DBConfig{
		host:                   configUtil.MustGetString("DB_HOST"),
		port:                   configUtil.MustGetString("DB_PORT"),
		name:                   configUtil.MustGetString("DB_NAME"),
		user:                   configUtil.MustGetString("DB_USER"),
		password:               configUtil.MustGetString("DB_PASSWORD"),
		sslmode:                configUtil.MustGetString("DB_SSL_MODE"),
		maxPoolSize:            configUtil.MustGetInt("DB_POOL_SIZE"),
		maxIdleConnections:     configUtil.MustGetInt("DB_MAX_IDLE_CONNECTIONS"),
		connMaxLifetimeMinutes: configUtil.MustGetInt("DB_CONN_MAX_LIFETIME_MINUTES"),
		transactionTimeoutInMS: configUtil.MustGetInt("DB_TRANSACTION_TIMEOUT_IN_MS"),
		isTest:                 configUtil.MustGetBool("IS_TEST"),
	}

	readDBCfg = DBConfig{
		host:                   getStringOrPanic("READ_DB_HOST"),
		port:                   getStringOrPanic("READ_DB_PORT"),
		name:                   getStringOrPanic("DB_NAME"),
		user:                   getStringOrPanic("DB_USER"),
		password:               getStringOrPanic("DB_PASSWORD"),
		sslmode:                getStringOrPanic("DB_SSL_MODE"),
		maxPoolSize:            getIntOrPanic("READ_DB_POOL_SIZE"),
		maxIdleConnections:     getIntOrPanic("READ_DB_MAX_IDLE_CONNECTIONS"),
		connMaxLifetimeMinutes: configUtil.MustGetInt("DB_CONN_MAX_LIFETIME_MINUTES"),
		transactionTimeoutInMS: getIntOrPanic("READ_DB_TRANSACTION_TIMEOUT_IN_MS"),
	}

	if dbCfg.isTest {
		dbCfg.name = configUtil.MustGetString("TEST_DB_NAME")
		dbCfg.password = configUtil.MustGetString("TEST_DB_PASSWORD")

		readDBCfg.name = configUtil.MustGetString("TEST_DB_NAME")
		readDBCfg.password = configUtil.MustGetString("TEST_DB_PASSWORD")
	}
}

func (config DBConfig) GetConnectionString() string {
	if config.isTest {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC password=%s", config.host, config.port, config.user, config.name, config.sslmode, config.password)

	}
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", config.host, config.port, config.user, config.name, config.sslmode, config.password)
}

func (config DBConfig) MaxPoolSize() int {
	return config.maxPoolSize
}

func (config DBConfig) ConnMaxLifetimeMinutes() int {
	return config.connMaxLifetimeMinutes
}

func (config DBConfig) Host() string {
	return config.host
}

func (config DBConfig) Name() string {
	return config.name
}

func (config DBConfig) MaxIdleConnections() int {
	return config.maxIdleConnections
}

func (config DBConfig) TransactionTimoutInMS() int {
	return config.transactionTimeoutInMS
}
