// Package secrets validates and parses all sensitive configuration.
package secrets

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/RMI/pacta/keyutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PACTAConfig struct {
	AuthVerificationKey AuthVerificationKey
	Postgres            *pgxpool.Config
}

type AuthVerificationKey struct {
	ID        string
	PublicKey ed25519.PublicKey
}

type RawPACTAConfig struct {
	PostgresConfig      *RawPostgresConfig
	AuthVerificationKey *RawAuthVerificationKey
}

type RawAuthVerificationKey struct {
	ID   string
	Data string
}

func LoadPACTA(rawCfg *RawPACTAConfig) (*PACTAConfig, error) {
	if rawCfg == nil {
		return nil, errors.New("no raw config provided")
	}

	pgxCfg, err := loadPGXConfig(rawCfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load 'postgres' config: %w", err)
	}

	authVerificationKey, err := parseAuthVerificationKey(rawCfg.AuthVerificationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse auth verification key config: %w", err)
	}

	return &PACTAConfig{
		Postgres:            pgxCfg,
		AuthVerificationKey: authVerificationKey,
	}, nil
}

type RawPostgresConfig struct {
	Host     string
	Port     int
	Database string

	User     string
	Password string
}

func loadPGXConfig(pgCfg *RawPostgresConfig) (*pgxpool.Config, error) {
	if pgCfg == nil {
		return nil, errors.New("config contained no 'postgres' configuration")
	}

	if pgCfg.Port < 0 {
		return nil, fmt.Errorf("invalid port %d given", pgCfg.Port)
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

func parseAuthVerificationKey(avk *RawAuthVerificationKey) (AuthVerificationKey, error) {
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
	in = strings.ReplaceAll(in, `\n`, "\n")
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
