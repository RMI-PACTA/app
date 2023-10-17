// Command server runs the PACTA API service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/RMI/pacta/azure/azblob"
	"github.com/RMI/pacta/azure/aztask"
	"github.com/RMI/pacta/cmd/runner/taskrunner"
	"github.com/RMI/pacta/cmd/server/pactasrv"
	"github.com/RMI/pacta/db/sqldb"
	"github.com/RMI/pacta/dockertask"
	"github.com/RMI/pacta/oapierr"
	oapipacta "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/secrets"
	"github.com/RMI/pacta/task"
	"github.com/Silicon-Ally/cryptorand"
	"github.com/Silicon-Ally/zaphttplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/namsral/flag"
	"github.com/rs/cors"
	"go.uber.org/zap"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("args cannot be empty")
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port = fs.Int("port", 8081, "Port for HTTP server")

		rateLimitMaxRequests = fs.Int("rate_limit_max_requests", 100, "The maximum number of requests to allow per rate_limit_unit_time before rate limiting the caller.")
		rateLimitUnitTime    = fs.Duration("rate_limit_unit_time", 1*time.Minute, "The unit of time over which to measure the rate_limit_max_requests.")

		allowedCORSOrigin = fs.String("allowed_cors_origin", "", "If specified, enables CORS handling and allows the given domain, e.g. 'http://localhost:3000'. This is used for the example web client in frontend/")

		env      = fs.String("env", "", "The environment that we're running in.")
		localDSN = fs.String("local_dsn", "", "If set, override the DB addresses retrieved from the secret configuration. Can only be used when running locally.")

		// Only when running locally because the Dockerized runner can't use local `az` CLI credentials
		localDockerTenantID     = fs.String("local_docker_tenant_id", "", "The Azure Tenant ID the localdocker service principal lives in")
		localDockerClientID     = fs.String("local_docker_client_id", "", "The client ID of the localdocker service principal")
		localDockerClientSecret = fs.String("local_docker_client_secret", "", "The client secret for accessing the localdocker service principal")

		// PACTA Execution
		useAZRunner = fs.Bool("use_azure_runner", false, "If true, execute PACTA on Azure Container Apps Jobs instead of a local instance.")

		// Secrets
		pgHost     = fs.String("secret_postgres_host", "", "Host of the Postgres server, like db.example.com")
		pgPort     = fs.Int("secret_postgres_port", 5432, "Port to connect to the Postgres server on")
		pgDatabase = fs.String("secret_postgres_database", "", "Name of the postgres database, like pactasrv")
		pgUser     = fs.String("secret_postgres_user", "", "Name of the Postgres user to connect as")
		pgPassword = fs.String("secret_postgres_password", "", "Password of the Postgres user to connect as")

		authKeyID   = fs.String("secret_auth_public_key_id", "", "Key ID (kid) of the JWT tokens to allow")
		authKeyData = fs.String("secret_auth_public_key_data", "", "PEM-encoded Ed25519 public key to verify JWT tokens with, contains literal \\n characters that will need to be replaced before parsing")

		azStorageAccount           = fs.String("secret_azure_storage_account", "", "The storage account to authenticate against for blob operations")
		azSourcePortfolioContainer = fs.String("secret_azure_source_portfolio_container", "", "The container in the storage account where we write raw portfolios to")

		runnerConfigLocation   = fs.String("secret_runner_config_location", "", "Location (like 'centralus') where the runner jobs should be executed")
		runnerConfigConfigPath = fs.String("secret_runner_config_config_path", "", "Config path (like '/configs/dev.conf') where the runner jobs should read their base config from")

		runnerConfigIdentityName               = fs.String("secret_runner_config_identity_name", "", "Name of the Azure identity to run runner jobs with")
		runnerConfigIdentitySubscriptionID     = fs.String("secret_runner_config_identity_subscription_id", "", "Subscription ID of the identity to run runner jobs with")
		runnerConfigIdentityResourceGroup      = fs.String("secret_runner_config_identity_resource_group", "", "Resource group of the identity to run runner jobs with")
		runnerConfigIdentityClientID           = fs.String("secret_runner_config_identity_client_id", "", "Client ID of the identity to run runner jobs with")
		runnerConfigIdentityManagedEnvironment = fs.String("secret_runner_config_identity_managed_environment", "", "Name of the Container Apps Environment where runner jobs should run")

		runnerConfigImageRegistry = fs.String("secret_runner_config_image_registry", "", "Registry where PACTA runner images live, like 'rmipacta.azurecr.io'")
		runnerConfigImageName     = fs.String("secret_runner_config_image_name", "", "Name of the Docker image of the PACTA runner, like 'runner'")
	)
	// Allows for passing in configuration via a -config path/to/env-file.conf
	// flag, see https://pkg.go.dev/github.com/namsral/flag#readme-usage
	fs.String(flag.DefaultConfigFlagname, "", "path to config file")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		logger *zap.Logger
		err    error
	)
	if *env == "local" {
		if logger, err = zap.NewDevelopment(); err != nil {
			return fmt.Errorf("failed to init logger: %w", err)
		}
	} else {
		if logger, err = zap.NewProduction(); err != nil {
			return fmt.Errorf("failed to init logger: %w", err)
		}
	}
	defer logger.Sync()

	sec, err := secrets.LoadPACTA(&secrets.RawPACTAConfig{
		PostgresConfig: &secrets.RawPostgresConfig{
			Host:     *pgHost,
			Port:     *pgPort,
			Database: *pgDatabase,
			User:     *pgUser,
			Password: *pgPassword,
		},
		AuthVerificationKey: &secrets.RawAuthVerificationKey{
			ID:   *authKeyID,
			Data: *authKeyData,
		},
		RunnerConfig: &secrets.RawRunnerConfig{
			Location:   *runnerConfigLocation,
			ConfigPath: *runnerConfigConfigPath,
			Identity: &secrets.RawRunnerIdentity{
				Name:               *runnerConfigIdentityName,
				SubscriptionID:     *runnerConfigIdentitySubscriptionID,
				ResourceGroup:      *runnerConfigIdentityResourceGroup,
				ClientID:           *runnerConfigIdentityClientID,
				ManagedEnvironment: *runnerConfigIdentityManagedEnvironment,
			},
			Image: &secrets.RawRunnerImage{
				Registry: *runnerConfigImageRegistry,
				Name:     *runnerConfigImageName,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to parse secrets: %w", err)
	}
	runCfg := sec.RunnerConfig

	if *localDSN != "" && *env != "local" {
		return errors.New("--local_dsn set outside of local environment")
	}

	var postgresCfg *pgxpool.Config
	if *localDSN != "" {
		if postgresCfg, err = pgxpool.ParseConfig(*localDSN); err != nil {
			return fmt.Errorf("failed to parse local DSN: %w", err)
		}
	} else {
		postgresCfg = sec.Postgres
	}

	logger.Info("Connecting to database", zap.String("db_host", postgresCfg.ConnConfig.Host))
	pgConn, err := pgxpool.NewWithConfig(ctx, postgresCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Pinging database")
	if err := pgConn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	db, err := sqldb.New(pgConn)
	if err != nil {
		return fmt.Errorf("failed to init sqldb: %w", err)
	}

	pactaSwagger, err := oapipacta.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to load PACTA swagger spec: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	pactaSwagger.Servers = nil

	var creds azcore.TokenCredential
	// This is necessary because the default timeout is too low in
	// azidentity.NewDefaultAzureCredentials, so it times out and fails to run.
	if azClientID := os.Getenv("AZURE_CLIENT_ID"); azClientID != "" {
		logger.Info("Loading user managed credentials", zap.String("client_id", azClientID))
		if creds, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(azClientID),
		}); err != nil {
			return fmt.Errorf("failed to load Azure credentials: %w", err)
		}
	} else {
		logger.Info("Loading default credentials")
		if creds, err = azidentity.NewDefaultAzureCredential(nil); err != nil {
			return fmt.Errorf("failed to load Azure credentials: %w", err)
		}
	}

	var runner taskrunner.Runner
	if *useAZRunner {
		logger.Info("initializing Azure task runner client")
		tmp, err := aztask.NewRunner(creds, &aztask.Config{
			Location: runCfg.Location,
			Rand:     rand.New(cryptorand.New()),
			Identity: &aztask.RunnerIdentity{
				Name:               runCfg.Identity.Name,
				SubscriptionID:     runCfg.Identity.SubscriptionID,
				ResourceGroup:      runCfg.Identity.ResourceGroup,
				ClientID:           runCfg.Identity.ClientID,
				ManagedEnvironment: runCfg.Identity.ManagedEnvironment,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to init Azure runner: %w", err)
		}
		runner = tmp
	} else {
		tmp, err := dockertask.NewRunner(logger, &dockertask.ServicePrincipal{
			TenantID:     *localDockerTenantID,
			ClientID:     *localDockerClientID,
			ClientSecret: *localDockerClientSecret,
		})
		if err != nil {
			return fmt.Errorf("failed to init docker runner: %w", err)
		}
		runner = tmp
	}

	tr, err := taskrunner.New(&taskrunner.Config{
		ConfigPath: runCfg.ConfigPath,
		BaseImage: &task.BaseImage{
			Registry: runCfg.Image.Registry,
			Name:     runCfg.Image.Name,
		},
		Logger: logger,
		Runner: runner,
	})
	if err != nil {
		return fmt.Errorf("failed to init task runner: %w", err)
	}

	blobClient, err := azblob.NewClient(creds, *azStorageAccount)
	if err != nil {
		return fmt.Errorf("failed to init blob client: %w", err)
	}

	// Create an instance of our handler which satisfies each generated interface
	srv := &pactasrv.Server{
		Blob:              blobClient,
		PorfolioUploadURI: *azSourcePortfolioContainer,
		DB:                db,
		TaskRunner:        tr,
	}

	pactaStrictHandler := oapipacta.NewStrictHandlerWithOptions(srv, nil /* middleware */, oapipacta.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: requestErrorHandlerFuncForService(logger, "pacta"),
		ResponseErrorHandlerFunc: oapierr.ErrorHandlerFunc(logger, func(err *oapierr.Error) *oapipacta.Error {
			// We don't care if it's the default message or not.
			cm, _ := err.ClientMessage()
			return &oapipacta.Error{
				ErrorId: string(err.ErrorID()),
				Message: cm,
			}
		}),
	})

	r := chi.NewRouter()

	jwKey, err := jwk.FromRaw(sec.AuthVerificationKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to make JWK key: %w", err)
	}
	jwKey.Set(jwk.KeyIDKey, sec.AuthVerificationKey.ID)

	// We now register our PACTA above as the handler for the interface
	oapipacta.HandlerWithOptions(pactaStrictHandler, oapipacta.ChiServerOptions{
		BaseRouter: r.With(
			// The order of these is important. We run RequestID and RealIP first to
			// populate relevant metadata for logging, and we run recovery immediately after
			// logging so it can catch any subsequent panics, but still has access to the
			// LogEntry created by the logging middleware.
			chimiddleware.RequestID,
			chimiddleware.RealIP,
			zaphttplog.NewMiddleware(logger),
			chimiddleware.Recoverer,

			jwtauth.Verifier(jwtauth.New("EdDSA", nil, jwKey)),
			jwtauth.Authenticator,

			oapimiddleware.OapiRequestValidator(pactaSwagger),

			rateLimitMiddleware(*rateLimitMaxRequests, *rateLimitUnitTime),
		),
	})

	// Created with https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Pagga&text=%20%20RMI%0APACTA
	fmt.Println()
	fmt.Println(`
     █▀▄░█▄█░▀█▀    
     █▀▄░█░█░░█░    
     ▀░▀░▀░▀░▀▀▀    
░█▀█░█▀█░█▀▀░▀█▀░█▀█
░█▀▀░█▀█░█░░░░█░░█▀█
░▀░░░▀░▀░▀▀▀░░▀░░▀░▀`)
	fmt.Println()

	// If CORS was specified, wrap our handler in that.
	var handler http.Handler
	if *allowedCORSOrigin != "" {
		handler = cors.New(cors.Options{
			AllowedOrigins:   []string{*allowedCORSOrigin},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Authorization", "Content-Type"},
			// Enable Debugging for testing, consider disabling in production
			Debug:          true,
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		}).Handler(r)
	} else {
		handler = r
	}

	s := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", *port),
	}

	// And we serve HTTP until the world ends.
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("error running HTTP server: %w", err)
	}

	return nil
}

func rateLimitMiddleware(maxReq int, windowLength time.Duration) func(http.Handler) http.Handler {
	// This example uses an in-memory rate limiter for simplicity, an application
	// that will be running multiple API instances should likely use something like
	// https://github.com/go-chi/httprate-redis to account for traffic across the
	// fleet.
	return httprate.Limit(
		maxReq,
		windowLength,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				return "", fmt.Errorf("failed to get claims from context: %w", err)
			}
			id, err := findFirstInClaims(claims, "user_id", "sub")
			if err != nil {
				return "", fmt.Errorf("failed to load user identifier: %w", err)
			}
			return id, nil
		}))
}

func findFirstInClaims(claims map[string]any, keys ...string) (string, error) {
	for _, k := range keys {
		v, ok := claims[k]
		if !ok {
			continue
		}
		vStr, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("%q claim was of unexpected type %T, wanted a string", k, v)
		}
		return vStr, nil
	}

	return "", errors.New("no valid claim was found")
}

func requestErrorHandlerFuncForService(logger *zap.Logger, svc string) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		// We log these at WARN because in aggregate, they might indicate an issue with our request handling.
		logger.Warn("error while parsing request", zap.String("service", svc), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func responseErrorHandlerFuncForService(logger *zap.Logger, svc string) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("error while handling request", zap.String("service", svc), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
