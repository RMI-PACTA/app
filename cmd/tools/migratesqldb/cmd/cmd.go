package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Silicon-Ally/testpgx/migrate"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Global objects
var (
	sqlDB *sql.DB
)

// Flags
var (
	dsn string // --dsn
)

// Commands
var (
	rootCmd = &cobra.Command{
		Use:   "migratesqldb",
		Short: "A simple tool for applying our migration set, using golang-migrate",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dsn == "" {
				return errors.New("no --dsn was specified")
			}

			pgCfg, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				return fmt.Errorf("failed to parse DSN: %w", err)
			}

			db, err := sql.Open("pgx", pgCfg.ConnString())
			if err != nil {
				return fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
			}
			sqlDB = db
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := sqlDB.Close(); err != nil {
				return fmt.Errorf("failed to close DB: %w", err)
			}
			return nil
		},
	}

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationsPath, err := bazel.Runfile("db/sqldb/migrations")
			if err != nil {
				return fmt.Errorf("failed to get a path to migrations: %w", err)
			}
			migrator, err := migrate.New(migrationsPath)
			if err != nil {
				return fmt.Errorf("when creating the migrator: %w", err)
			}
			if err := migrator.Migrate(sqlDB); err != nil {
				return fmt.Errorf("while applying the migration(s): %w", err)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "A Postgres DSN, parsable by pgx.ParseConfig")
	rootCmd.AddCommand(applyCmd)
}
