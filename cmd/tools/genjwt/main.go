// Command genjwt creates and signs a JWT token that can be used with the PACTA
// API.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/RMI/pacta/keyutil"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		keyType = flag.String("key_type", "ed25519", "The type of the key pair to load, currently only 'ed25519' is supported")

		privKeyFile = flag.String("private_key_file", "test_server.key", "The path to read the private key to use for signing from.")
		userID      = flag.String("user_id", "test123", "The ID of the user to put in the 'sub' claim of the token.")
		expiresIn   = flag.String("expires_in", "24h", "When the token should expire, relative to now. Should be formatted in a way that time.ParseDuration can handle.")
	)
	flag.Parse()

	switch *keyType {
	case "ed25519":
		// This is the only one we currently support.
	default:
		return fmt.Errorf("unsupported key type %q", *keyType)
	}

	expiresInDur, err := time.ParseDuration(*expiresIn)
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}

	priv, err := keyutil.DecodeED25519PrivateKeyFromFile(*privKeyFile)
	if err != nil {
		return fmt.Errorf("failed to load private key file: %w", err)
	}

	jwtAuth := jwtauth.New("EdDSA", priv, nil /* verify, unused */)
	now := time.Now()
	claims := map[string]any{
		"sub": *userID,
		"exp": now.Add(expiresInDur),
		"nbf": now.Add(-5 * time.Second),
	}

	_, tkn, err := jwtAuth.Encode(claims)
	if err != nil {
		return fmt.Errorf("failed to sign token: %w", err)
	}

	fmt.Printf("\n\nToken: %s\n\n", tkn)

	return nil
}
