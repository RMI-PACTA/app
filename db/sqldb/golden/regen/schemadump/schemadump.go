package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/testpgx"
	"github.com/Silicon-Ally/testpgx/migrate"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	migrationsPath, err := bazel.Runfile("db/sqldb/migrations")
	if err != nil {
		log.Fatalf("failed to get a path to migrations: %v", err)
	}
	migrator, err := migrate.New(migrationsPath)
	if err != nil {
		log.Fatalf("failed to init migrator: %v", err)
	}
	ctx := context.Background()

	env, err := testpgx.New(ctx, testpgx.WithMigrator(migrator), testpgx.WithPostgresDockerImage("postgres:14.9"))
	if err != nil {
		log.Fatalf("while creating/getting the test env: %v", err)
	}
	defer func() {
		err = env.TearDown(ctx)
		if err != nil {
			log.Fatalf("while trying to tear down env: %v", err)
		}
	}()

	env.WithMigratedDB(ctx, func(db *pgxpool.Pool) error {
		result, err := env.DumpDatabaseSchema(ctx, db.Config().ConnConfig.Database)
		if err != nil {
			log.Fatalf("while dumping database schema: %v", err)
		}
		fmt.Printf(result)
		return nil
	})
	os.Exit(0)
}
