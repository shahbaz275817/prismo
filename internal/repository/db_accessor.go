package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
	"github.com/shahbaz275817/prismo/pkg/reporting"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	POSTGRES                    = "postgres"
	txKey                   key = "dbtx"
	txTimeoutErrorCode          = "57014"
	dbMetricKey                 = "database"
	dbMetricPublishDuration     = 10 * time.Second
)

type key string

type NewRelicDetail struct {
	Operation string
	Target    string
}

type Accessor interface {
	Close() error
	Ping() error
	Transact(context.Context, func(ctx context.Context) error) error
	TransactWithTimeout(context.Context, time.Duration, func(ctx context.Context) error, *sql.TxOptions) error
	TransactWithoutTx(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) error) error
}

type accessor struct {
	Db          *gorm.DB
	txTimeoutMS int
}

func (acc accessor) Ping() error {
	sqlDbInstance, err := acc.Db.DB()
	if err != nil {
		logger.Errorf("Failed to get DB instance: %s", err)
		return err
	}
	return sqlDbInstance.Ping()
}

func (acc accessor) Close() error {
	logger.WithContext(context.Background()).Infof("Closing dB conn pool")
	sqlDbInstance, err := acc.Db.DB()
	if err != nil {
		logger.Errorf("Failed to get DB instance: %s", err)
		return err
	}
	return sqlDbInstance.Close()
}

func NewAccessor(dbConfig config.DBConfig) (Accessor, error) {
	logger.Infof("Creating dB conn pool: %s %s", dbConfig.Host(), dbConfig.Name())

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             500 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  gormLogger.Info,        // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                  // Don't include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbConfig.GetConnectionString()), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger:      newLogger,
		PrepareStmt: false,
	})
	if err != nil {
		logger.WithContext(context.Background()).Error("Failed to load Database")
		return nil, err
	}
	sqlDbInstance, err := db.DB()
	if err = sqlDbInstance.Ping(); err != nil {
		logger.WithContext(context.Background()).Errorf("ping to the database host failed: %s", err)
		return nil, err
	}

	sqlDbInstance.SetMaxIdleConns(dbConfig.MaxIdleConnections())
	sqlDbInstance.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetimeMinutes()) * time.Minute)
	sqlDbInstance.SetMaxOpenConns(dbConfig.MaxPoolSize())

	logger.Infof("Created dB conn pool")

	return accessor{db, dbConfig.TransactionTimoutInMS()}, nil
}

func contextWithDBTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		logger.WithContext(ctx).Error("No DB transaction found in context")
		return nil
	}
	return tx
}

func (acc accessor) TransactWithoutTx(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) error) error {

	return fn(ctx, acc.Db)
}

func (acc accessor) Transact(ctx context.Context, txFunc func(context.Context) error) (err error) {
	return acc.TransactWithTimeout(ctx, time.Duration(acc.txTimeoutMS*int(time.Millisecond)), txFunc, &sql.TxOptions{})
}

func (acc accessor) TransactWithTimeout(ctx context.Context, timeout time.Duration, txFunc func(context.Context) error, txOpts *sql.TxOptions) (err error) {

	isOwner := false

	tx, txExists := ctx.Value(txKey).(*gorm.DB)

	if !txExists {
		isOwner = true

		ctx, cancelFunc := context.WithTimeout(ctx, timeout)
		defer cancelFunc()
		tx = acc.Db.WithContext(ctx).Begin(txOpts)
		if tx.Error != nil {
			return tx.Error
		}
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = errors.WithStack(p)
			default:
				err = errors.Errorf("%s", p)
			}
		}
		if err != nil {
			if isOwner {
				tx.Rollback()
			} else {
				err = errors.WithStack(err)
				//panic(err) // pass thorough panic to the outermost tx --- WHY?
			}
			return
		}
		if isOwner { // commit only if outermost transaction
			tx.Commit()
		}
	}()

	err = txFunc(contextWithDBTx(ctx, tx))
	if err != nil {
		logger.Infof("Error in DB Transaction: %s", err.Error())
	}
	return err
}

func reportDBMetric(err error, entry *reporting.ReporterEntry) {
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if checkDBTimeoutError(err) {
		entry.Timeout()
		return
	}

	entry.Failure()
}

func checkDBTimeoutError(err error) bool {
	if err == context.DeadlineExceeded {
		return true
	}

	if pgconn.Timeout(err) {
		return true
	}

	return false
}
