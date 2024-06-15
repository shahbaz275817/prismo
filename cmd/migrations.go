package main

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // get db migration from path
	_ "github.com/lib/pq"

	"github.com/shahbaz275817/prismo/pkg/errors"
)

const (
	defaultMigrationsPath = "file://./migrations"
	postgresDriver        = "postgres"
)

var ErrNoMigrations = errors.New("no migrations")

func RunDatabaseMigrations(connectionURL, path string) error {
	m, err := createMigration(connectionURL, path)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange || err == nil {
		return nil
	}

	return err
}

func RollbackLatestMigration(connectionURL, path string) error {
	m, err := createMigration(connectionURL, path)
	if err != nil {
		return err
	}

	err = m.Steps(-1)
	if err == migrate.ErrNoChange || err == nil {
		return nil
	}

	return err
}

func createMigration(connectionURL, path string) (*migrate.Migrate, error) {
	if path == "" {
		path = defaultMigrationsPath
	}

	db, err := sql.Open(postgresDriver, connectionURL)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(path, postgresDriver, driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}
