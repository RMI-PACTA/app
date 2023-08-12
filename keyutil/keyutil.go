// Package keyutil provides some simple wrappers for serializing + deserializing
// cryptographic keys.
package keyutil

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func GenerateED25519ToFiles(pubFile, privFile string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ED25519 key: %w", err)
	}

	if err := EncodeED25519PublicKeyToFile(pub, pubFile); err != nil {
		return fmt.Errorf("failed to encode public key to file: %w", err)
	}
	if err := EncodeED25519PrivateKeyToFile(priv, privFile); err != nil {
		return fmt.Errorf("failed to encode private key to file: %w", err)
	}

	return nil
}

func EncodeED25519PublicKeyToFile(pub ed25519.PublicKey, out string) error {
	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	if err := EncodeToFile(out, "PUBLIC KEY", pubDER); err != nil {
		return fmt.Errorf("failed to encode and write public key: %w", err)
	}

	return nil
}

func EncodeED25519PrivateKeyToFile(priv ed25519.PrivateKey, out string) error {
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := EncodeToFile(out, "PRIVATE KEY", privDER); err != nil {
		return fmt.Errorf("failed to encode and write private key: %w", err)
	}

	return nil
}

func EncodeToFile(name, typ string, dat []byte) error {
	block := &pem.Block{
		Type:  typ,
		Bytes: dat,
	}

	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := pem.Encode(f, block); err != nil {
		return fmt.Errorf("failed to PEM encode data: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}

func DecodeED25519PublicKeyFromFile(in string) (ed25519.PublicKey, error) {
	pubDER, err := DecodeFromFile(in, "PUBLIC KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to PEM decode public key file: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(pubDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data as PKIX ASN.1 DER-formatted public key: %w", err)
	}

	pubED, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key was of type %T, expected ed25519.PublicKey", pub)
	}

	return pubED, nil
}

func DecodeED25519PrivateKeyFromFile(in string) (ed25519.PrivateKey, error) {
	privDER, err := DecodeFromFile(in, "PRIVATE KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to PEM decode private key file: %w", err)
	}

	priv, err := x509.ParsePKCS8PrivateKey(privDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data as PKIX ASN.1 DER-formatted private key: %w", err)
	}

	privED, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key was of type %T, expected ed25519.PrivateKey", priv)
	}

	return privED, nil
}

func DecodeFromFile(name, typ string) ([]byte, error) {
	dat, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM-encoded file: %w", err)
	}

	block, _ := pem.Decode(dat)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	if block.Type != typ {
		return nil, fmt.Errorf("block type was %q, expected %q", block.Type, typ)
	}

	return block.Bytes, nil
}
