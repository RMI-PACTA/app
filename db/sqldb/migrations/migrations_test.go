package migrations

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Silicon-Ally/testpgx"
	"github.com/Silicon-Ally/testpgx/migrate"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
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
		log.Fatalf(" while creating/getting the test env: %v", err)
	}
	defer func() {
		err = env.TearDown(ctx)
		if err != nil {
			log.Fatalf("while trying to tear down env: %v", err)
		}
	}()
	sqlEnv = env
	result := m.Run()
	time.Sleep(10 * time.Second)
	return result
}

var sqlEnv *testpgx.Env

type migrationFile struct {
	num  int64
	name string
	up   bool
}

func TestListFiles(t *testing.T) {
	re := regexp.MustCompile(`(\d+)_(.+)\.(up|down)\.sql`)

	rfp, err := bazel.RunfilesPath()
	if err != nil {
		t.Fatalf("when getting runfiles path: %v", err)
	}
	p := path.Join(rfp, "db", "sqldb", "migrations")
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		t.Fatalf("couldn't read from migrations dir: %v", err)
	}
	var files []migrationFile
	for _, fi := range fis {
		n := fi.Name()
		if !strings.HasSuffix(n, ".sql") {
			continue
		}
		match := re.FindStringSubmatch(n)
		if match == nil {
			t.Fatalf("file %q did not conform to regex expectations", n)
		}
		i, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			t.Fatalf("couldn't parse int %q", match[1])
		}
		files = append(files, migrationFile{
			num:  i,
			name: match[2],
			up:   match[3] == "up",
		})
	}
	if len(files) == 0 {
		t.Fatalf("no matching files found at path %q", p)
	}

	m := make(map[int64]map[string]map[bool]bool)
	for _, f := range files {
		if _, ok := m[f.num]; !ok {
			m[f.num] = make(map[string]map[bool]bool)
		}
		if _, ok := m[f.num][f.name]; !ok {
			m[f.num][f.name] = make(map[bool]bool)
		}
		if m[f.num][f.name][f.up] {
			t.Fatalf("collision for value %+v", f)
		}
		m[f.num][f.name][f.up] = true
	}

	for num, names := range m {
		if len(names) != 1 {
			t.Fatalf("multiple migrations sharing the same id %d: %+v", num, names)
		}
		for name, updowns := range names {
			if len(updowns) != 2 {
				t.Fatalf("expected an up and down migration for value %q, but got %+v", name, updowns)
			}
		}
	}
}

func TestRegeneratedHumanReadableSchema(t *testing.T) {
	rfp, err := bazel.RunfilesPath()
	if err != nil {
		t.Fatalf("when getting runfiles path: %v", err)
	}
	p := path.Join(rfp, "db/sqldb/golden/human_readable_schema.sql")
	want, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("couldn't read human schema file: %v", err)
	}
	ctx := context.Background()
	db := sqlEnv.GetMigratedDB(ctx, t)
	got, err := sqlEnv.DumpDatabaseSchema(ctx, db.Config().ConnConfig.Database, testpgx.WithHumanReadableSchema())
	if err != nil {
		t.Fatalf("failed to dump database schema: %v", err)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("unexpected diff in human readable schema: \n%s", diff)
		fmt.Fprint(os.Stderr, "\nTo fix, run the following:\n")
		fmt.Fprint(os.Stderr, "\t$ bazel run //scripts:regen_db_goldens\n\n")
	}
}

func TestRegeneratedDatabaseSchemaDump(t *testing.T) {
	rfp, err := bazel.RunfilesPath()
	if err != nil {
		t.Fatalf("when getting runfiles path: %v", err)
	}
	p := path.Join(rfp, "db/sqldb/golden/schema_dump.sql")
	want, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("couldn't read database schema file: %v", err)
	}
	ctx := context.Background()
	db := sqlEnv.GetMigratedDB(ctx, t)
	got, err := sqlEnv.DumpDatabaseSchema(ctx, db.Config().ConnConfig.Database)
	if err != nil {
		t.Fatalf("failed to dump database schema: %v", err)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("unexpected diff in datbase schema dump: \n%s", diff)
		fmt.Fprint(os.Stderr, "\nTo fix, run the following:\n")
		fmt.Fprint(os.Stderr, "\t$ bazel run //scripts:regen_db_goldens\n\n")
	}
}
