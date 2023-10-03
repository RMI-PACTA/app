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
	RunnerConfig        RunnerConfig
}

type AuthVerificationKey struct {
	ID        string
	PublicKey ed25519.PublicKey
}

type RunnerConfig struct {
	Location   string
	ConfigPath string
	Identity   RunnerIdentity
	Image      RunnerImage
}

type RunnerIdentity struct {
	Name               string
	SubscriptionID     string
	ResourceGroup      string
	ClientID           string
	ManagedEnvironment string
}

type RunnerImage struct {
	Registry string
	Name     string
}

type RawPACTAConfig struct {
	PostgresConfig      *RawPostgresConfig
	AuthVerificationKey *RawAuthVerificationKey
	RunnerConfig        *RawRunnerConfig
}

type RawAuthVerificationKey struct {
	ID   string
	Data string
}

type RawRunnerConfig struct {
	Location   string
	ConfigPath string
	Identity   *RawRunnerIdentity
	Image      *RawRunnerImage
}

type RawRunnerIdentity struct {
	Name               string
	SubscriptionID     string
	ResourceGroup      string
	ClientID           string
	ManagedEnvironment string
}

type RawRunnerImage struct {
	Registry string
	Name     string
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

	runnerConfig, err := parseRunnerConfig(rawCfg.RunnerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse runner config: %w", err)
	}

	return &PACTAConfig{
		Postgres:            pgxCfg,
		AuthVerificationKey: authVerificationKey,
		RunnerConfig:        runnerConfig,
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

func parseRunnerConfig(cfg *RawRunnerConfig) (RunnerConfig, error) {
	if cfg == nil {
		return RunnerConfig{}, errors.New("no runner config was provided")
	}

	if cfg.Location == "" {
		return RunnerConfig{}, errors.New("no runner_config.location was provided")
	}

	if cfg.Image == nil {
		return RunnerConfig{}, errors.New("no runner_config.image was provided")
	}
	if cfg.Image.Name == "" {
		return RunnerConfig{}, errors.New("no runner_config.image.name was provided")
	}
	if cfg.Image.Registry == "" {
		return RunnerConfig{}, errors.New("no runner_config.image.registry was provided")
	}

	if cfg.ConfigPath == "" {
		return RunnerConfig{}, errors.New("no runner_config.config_path was provided")
	}

	if cfg.Identity == nil {
		return RunnerConfig{}, errors.New("no runner_config.identity was provided")
	}
	if cfg.Identity.Name == "" {
		return RunnerConfig{}, errors.New("no runner_config.identity.name was provided")
	}
	if cfg.Identity.SubscriptionID == "" {
		return RunnerConfig{}, errors.New("no runner_config.identity.subscription_id was provided")
	}
	if cfg.Identity.ResourceGroup == "" {
		return RunnerConfig{}, errors.New("no runner_config.identity.resource_group was provided")
	}
	if cfg.Identity.ClientID == "" {
		return RunnerConfig{}, errors.New("no runner_config.identity.client_id was provided")
	}
	if cfg.Identity.ManagedEnvironment == "" {
		return RunnerConfig{}, errors.New("no runner_config.identity.managed_environment was provided")
	}

	return RunnerConfig{
		Location:   cfg.Location,
		ConfigPath: cfg.ConfigPath,
		Identity: RunnerIdentity{
			Name:               cfg.Identity.Name,
			SubscriptionID:     cfg.Identity.SubscriptionID,
			ResourceGroup:      cfg.Identity.ResourceGroup,
			ClientID:           cfg.Identity.ClientID,
			ManagedEnvironment: cfg.Identity.ManagedEnvironment,
		},
		Image: RunnerImage{
			Registry: cfg.Image.Registry,
			Name:     cfg.Image.Name,
		},
	}, nil
}
