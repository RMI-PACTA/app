// Command server runs the PACTA API service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RMI/pacta/cmd/server/pactasrv"
	"github.com/RMI/pacta/db/sqldb"
	"github.com/RMI/pacta/keyutil"
	"github.com/RMI/pacta/oapierr"
	oapipacta "github.com/RMI/pacta/openapi/pacta"
	"github.com/Silicon-Ally/zaphttplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		localDSN = fs.String("local_dsn", "", "If set, override the DB addresses retrieved from the sops configuration. Can only be used when running locally.")

		authPubKeyFile = fs.String("auth_public_key_file", "", "The PEM-encoded PKIX ASN.1 DER-formatted ED25519 public key to verify JWTs with")
	)
	// Allows for passing in configuration via a -config path/to/env-file.conf
	// flag, see https://pkg.go.dev/github.com/namsral/flag#readme-usage
	fs.String(flag.DefaultConfigFlagname, "", "path to config file")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pub is the key we use to authenticate signatures on user auth tokens.
	pub, err := keyutil.DecodeED25519PublicKeyFromFile(*authPubKeyFile)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	var logger *zap.Logger
	if *env == "local" {
		if logger, err = zap.NewDevelopment(); err != nil {
			return fmt.Errorf("failed to init logger: %w", err)
		}
	} else {
		if logger, err = zap.NewProduction(); err != nil {
			return fmt.Errorf("failed to init logger: %w", err)
		}
	}

	if *localDSN != "" && *env != "local" {
		return errors.New("--local_dsn set outside of local environment")
	}

	var postgresCfg *pgxpool.Config
	if *localDSN != "" {
		if postgresCfg, err = pgxpool.ParseConfig(*localDSN); err != nil {
			return fmt.Errorf("failed to parse local DSN: %w", err)
		}
	} else {
		// TODO: Add support for sops-encrypted credentials.
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

	// Create an instance of our handler which satisfies each generated interface
	srv := &pactasrv.Server{
		DB: db,
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

			jwtauth.Verifier(jwtauth.New("EdDSA", nil, pub)),
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
