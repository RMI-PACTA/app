package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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
			migrationsPath, err := runfiles.Rlocation("__main__/db/sqldb/migrations/0001_create_schema_migrations_history.down.sql")
			if err != nil {
				return fmt.Errorf("failed to get a path to migrations: %w", err)
			}
			migrationsPath = filepath.Dir(migrationsPath)
			migrator, err := newMigrator(sqlDB, migrationsPath)
			if err != nil {
				return fmt.Errorf("when creating the migrator: %w", err)
			}
			if err := migrator.Up(); err != nil {
				return fmt.Errorf("while applying the migration(s): %w", err)
			}
			return nil
		},
	}

	rollbackCmd = &cobra.Command{
		Use:   "rollback",
		Short: "Rollback migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationsPath, err := runfiles.Rlocation("__main__/db/sqldb/migrations/0001_create_schema_migrations_history.down.sql")
			if err != nil {
				return fmt.Errorf("failed to get a path to migrations: %w", err)
			}
			migrationsPath = filepath.Dir(migrationsPath)
			migrator, err := newMigrator(sqlDB, migrationsPath)
			if err != nil {
				return fmt.Errorf("when creating the migrator: %w", err)
			}
			if err := migrator.Steps(-1); err != nil {
				return fmt.Errorf("while rolling back the migration(s): %w", err)

			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "A Postgres DSN, parsable by pgx.ParseConfig")
	rootCmd.AddCommand(applyCmd, rollbackCmd)
}

func newMigrator(db *sql.DB, migrationsPath string) (*migrate.Migrate, error) {
	// Pings the database to distinguish between migration and connection errors
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to init driver instance: %w", err)
	}

	mgr, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init migrate instance: %w", err)
	}

	return mgr, nil
}
