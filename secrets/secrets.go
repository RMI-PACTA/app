// Package secrets implements a wrapper around sops
// (https://github.com/mozilla/sops) that decrypts encrypted secret files on
// disk.
package secrets

import (
	"crypto/ed25519"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RMI/pacta/keyutil"
	"github.com/getsops/sops/v3/decrypt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// sopsFileExtension is the suffix we expect all
	// sops-encrypted configurations to have.
	sopsFileExtension = ".enc.json"
)

type PACTASecretsConfig struct {
	AuthVerificationKey AuthVerificationKey
	Postgres            *pgxpool.Config
}

type AuthVerificationKey struct {
	ID        string
	PublicKey ed25519.PublicKey
}

type pactaConfig struct {
	PostgresConfig      *postgresConfig      `json:"postgres"`
	AuthVerificationKey *authVerificationKey `json:"auth_public_key"`
}

type authVerificationKey struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func LoadPACTA(name string) (*PACTASecretsConfig, error) {
	var cfg pactaConfig
	if err := loadConfig(name, &cfg); err != nil {
		return nil, err
	}

	pgxCfg, err := loadPGXConfig(cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load 'postgres' config: %w", err)
	}

	authVerificationKey, err := parseAuthVerificationKey(cfg.AuthVerificationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse auth verification key config: %w", err)
	}

	return &PACTASecretsConfig{
		Postgres:            pgxCfg,
		AuthVerificationKey: authVerificationKey,
	}, nil
}

type MigratorConfig struct {
	Postgres *pgxpool.Config
}

type migratorConfig struct {
	PostgresConfig *postgresConfig `json:"postgres"`
}

type postgresConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Database string `json:"database"`

	User     string `json:"user"`
	Password string `json:"password"`
}

func LoadMigratorConfig(name string) (*MigratorConfig, error) {
	var cfg migratorConfig
	if err := loadConfig(name, &cfg); err != nil {
		return nil, err
	}

	pgxCfg, err := loadPGXConfig(cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load 'postgres' config: %w", err)
	}

	return &MigratorConfig{
		Postgres: pgxCfg,
	}, nil
}

func loadPGXConfig(pgCfg *postgresConfig) (*pgxpool.Config, error) {
	if pgCfg == nil {
		return nil, errors.New("config contained no 'postgres' configuration")
	}

	var kvs strings.Builder
	add := func(key, val string) {
		if val == "" {
			return
		}

		// Add a space to start if we aren't the first.
		if kvs.Len() > 0 {
			kvs.WriteRune(' ')
		}

		kvs.WriteString(key)
		kvs.WriteRune('=')
		kvs.WriteString(val)
	}

	add("host", pgCfg.Host)
	add("dbname", pgCfg.Database)
	add("user", pgCfg.User)
	add("password", pgCfg.Password)

	if pgCfg.Port > 0 {
		add("port", strconv.Itoa(int(pgCfg.Port)))
	}

	pgxCfg, err := pgxpool.ParseConfig(kvs.String())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig failed on kv pairs: %w", err)
	}

	return pgxCfg, nil
}

func loadConfig(name string, v interface{}) error {
	if err := checkFilename(name); err != nil {
		return err
	}

	dat, err := decrypt.File(name, "json")
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}

	if err := json.Unmarshal(dat, v); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func checkFilename(name string) error {
	fn := filepath.Base(name)
	if !strings.HasSuffix(fn, sopsFileExtension) {
		return fmt.Errorf("the given sops config %q does not have the expected extension %q", fn, sopsFileExtension)
	}
	return nil
}

func parseAuthVerificationKey(avk *authVerificationKey) (AuthVerificationKey, error) {
	if avk == nil {
		return AuthVerificationKey{}, errors.New("no auth_public_key was provided")
	}

	if avk.ID == "" {
		return AuthVerificationKey{}, errors.New("no auth_public_key.id was provided")
	}

	if avk.Data == "" {
		return AuthVerificationKey{}, errors.New("no auth_public_key.data was provided, should be PEM-encoded PKCS #8 ASN.1 DER-formatted ED25519 public key")
	}

	pub, err := loadPublicKey(avk.Data)
	if err != nil {
		return AuthVerificationKey{}, fmt.Errorf("failed to load auth verification key: %w", err)
	}
	return AuthVerificationKey{
		ID:        avk.ID,
		PublicKey: pub,
	}, nil
}

func loadPublicKey(in string) (ed25519.PublicKey, error) {
	pubDER, err := decodePEM("PUBLIC KEY", []byte(in))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PEM-encoded public key: %w", err)
	}

	pub, err := keyutil.DecodeED25519PublicKey(pubDER)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	return pub, nil
}

func decodePEM(typ string, dat []byte) ([]byte, error) {
	block, _ := pem.Decode(dat)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	if block.Type != typ {
		return nil, fmt.Errorf("block type was %q, expected %q", block.Type, typ)
	}

	return block.Bytes, nil
}
