package secrets

import (
	"testing"

	"github.com/Silicon-Ally/testsops"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCheckFilename(t *testing.T) {
	tests := []struct {
		desc    string
		path    string
		wantErr bool
	}{
		{
			desc: "valid relative path",
			path: "relative/path/to/dev.enc.json",
		},
		{
			desc: "valid absolute path",
			path: "/configs/secrets/local.enc.json",
		},
		{
			desc:    "invalid config",
			path:    "/configs/local.conf",
			wantErr: true,
		},
		{
			desc:    "invalid non-encrypted ext",
			path:    "/configs/secrets/prod.json",
			wantErr: true,
		},
		{
			desc:    "invalid blank path",
			path:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			gotErr := checkFilename(test.path)
			if test.wantErr && gotErr == nil {
				t.Error("no error was returned, but one was expected")
			}
			if !test.wantErr && gotErr != nil {
				t.Errorf("checkFilename(%q): %v", test.path, gotErr)
			}
		})
	}
}

func TestLoadMigratorConfig(t *testing.T) {
	tests := []struct {
		desc     string
		contents string
		want     *MigratorConfig
	}{
		{
			desc: "all parameters",
			contents: `{
  "postgres": {
    "host": "test-host",
    "port": 5434,
    "database": "db-name",
    "user": "postgres",
    "password": "not-a-real-password"
  }
}`,
			want: &MigratorConfig{
				Postgres: &pgxpool.Config{
					ConnConfig: &pgx.ConnConfig{
						Config: pgconn.Config{
							Host:     "test-host",
							Port:     5434,
							Database: "db-name",
							User:     "postgres",
							Password: "not-a-real-password",
						},
					},
				},
			},
		},
		{
			desc: "valid config with some parameters",
			contents: `{
  "postgres": {
    "host": "another-test-host",
    "database": "another-db-name",
    "user": "postgres"
  }
}`,
			want: &MigratorConfig{
				Postgres: &pgxpool.Config{
					ConnConfig: &pgx.ConnConfig{
						Config: pgconn.Config{
							Host:     "another-test-host",
							Database: "another-db-name",
							User:     "postgres",
							Port:     5432, // the default
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			sopsPath, ok := bazel.FindBinary("cmd/sops", "sops")
			if !ok {
				t.Fatal("'sops' binary not found in runfiles")
			}
			sopsCfg := testsops.EncryptJSON(t, test.contents, testsops.WithSOPSBinary(sopsPath))

			t.Setenv("SOPS_AGE_KEY_FILE", sopsCfg.KeyPath)
			got, err := LoadMigratorConfig(sopsCfg.EncryptedContentsPath)
			if err != nil {
				t.Fatalf("failed to load migrator config: %v", err)
			}

			if diff := cmp.Diff(test.want, got, compareMigratorConfigs()); diff != "" {
				t.Errorf("unexpected migrator config (-want +got)\n%s", diff)
			}
		})
	}
}

func compareMigratorConfigs() cmp.Option {
	// We add a custom comparer for *pgx.ConnConfig because it contains lots of
	// fields we don't actually care about.
	return cmp.Comparer(func(aCfg, bCfg *pgxpool.Config) bool {
		a := aCfg.ConnConfig.Config
		b := bCfg.ConnConfig.Config
		return a.Host == b.Host &&
			a.Port == b.Port &&
			a.Database == b.Database &&
			a.User == b.User &&
			a.Password == b.Password
	})
}
