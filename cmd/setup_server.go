package main

import (
	"context"
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/shahbaz275817/prismo/internal/repository/account"
	"github.com/shahbaz275817/prismo/internal/repository/operationtype"
	"github.com/shahbaz275817/prismo/internal/repository/transaction"
	account2 "github.com/shahbaz275817/prismo/internal/services/account"
	operationType2 "github.com/shahbaz275817/prismo/internal/services/operationtype"
	transaction2 "github.com/shahbaz275817/prismo/internal/services/transaction"
	"github.com/shahbaz275817/prismo/pkg/cache"
	"github.com/shahbaz275817/prismo/pkg/locks"

	appcontext "github.com/shahbaz275817/prismo/internal/appcontext/server"
	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/pkg/logger"
	"github.com/shahbaz275817/prismo/pkg/reporting"
)

func InitializeHandlerDependencies() (appcontext.Dependencies, func(), error) {
	return InitializeHandlerDependenciesWithExternal(appcontext.ExternalDependencies{})
}

func InitializeHandlerDependenciesWithExternal(ext appcontext.ExternalDependencies) (appcontext.Dependencies, func(), error) {

	db, err := repository.NewAccessor(config.DB())
	if err != nil {
		logger.Fatalf("unable to setup db: %v", err)
		return appcontext.Dependencies{}, nil, err
	}

	cacheClient, err := cache.NewClient(config.Cache())
	if err != nil {
		logger.Fatalf("unable to connect to redis: %v", err)
		return appcontext.Dependencies{}, nil, err
	}
	logger.Infof("Connection to Redis Cache success")
	atomicLock := locks.NewAtomicLock(cacheClient, config.AtomicLockConfig())

	accountRepository := account.NewAccountRepository(db)
	accountService := account2.NewAccountService(accountRepository)

	operationTypeRepository := operationtype.NewOperationTypeRepository(db)
	operationTypeService := operationType2.NewOperationTypeService(operationTypeRepository)

	transactionRepository := transaction.NewTransactionRepository(db)
	transactionService := transaction2.NewTransactionService(transactionRepository)

	return appcontext.Dependencies{
			AccountService:        accountService,
			OperationTypesService: operationTypeService,
			TransactionService:    transactionService,
			AtomicLock:            atomicLock,
		}, func() {
			db.Close()
		}, nil
}

func initReporter(cfg reporting.StatsDConfig) (reporter *reporting.Reporter, err error) {
	logger.Infof("Creating StatsD connection with %v %s %s", cfg.Enabled, cfg.Host, cfg.Namespace)

	metricReporter, err := reporting.NewStatsD(cfg)
	if err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"during": "statsd_initialization"})
		errDesc := fmt.Sprintf("failed to initialize StatsDReporter: %s", err)
		logger.WithContext(context.Background()).Errorf(errDesc)
		return nil, err
	}
	reporting.NewHystrixStatsDCollector(cfg)

	logger.Infof("Created StatsD connection")

	return &reporting.Reporter{
		MetricReporter: metricReporter,
	}, nil
}
