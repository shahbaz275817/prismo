package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/internal/handler"
	"github.com/shahbaz275817/prismo/pkg/logger"
	"github.com/shahbaz275817/prismo/pkg/server"
)

func newCLI() *cobra.Command {
	cli := &cobra.Command{
		Use:   "amphibian",
		Short: "amphibian is to store and maintain partner data",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			config.Load()
			logger.SetupLogger(config.Log())

		},
	}

	cli.AddCommand(newServerCmd())
	cli.AddCommand(newMigrateCmd())
	cli.AddCommand(newRollbackCmd())
	return cli
}

func newServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "server",
		Short:   "Start HTTP API server",
		Aliases: []string{"start"},
		Run: func(_ *cobra.Command, _ []string) {
			deps, cleanup, err := InitializeHandlerDependencies()
			if err != nil {
				logger.WithContext(context.Background()).Errorf("Server Dependencies initialization failed: %s", err)
				panic(err)
			}
			defer cleanup()

			s := server.New(handler.NewRouter(deps))
			s.Serve(config.Addr())
		},
	}
}

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Perform db migration",
		Run: func(_ *cobra.Command, _ []string) {
			if err := RunDatabaseMigrations(config.DB().GetConnectionString(), ""); err != nil {
				logger.Fatalf("Migrate: unable to run migration %v", err)
			}
		},
	}
}

func newRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback last db migration",
		Run: func(_ *cobra.Command, _ []string) {
			if err := RollbackLatestMigration(config.DB().GetConnectionString(), ""); err != nil {
				logger.Fatalf("Migrate: unable to rollback migration %v", err)
			}
		},
	}
}
