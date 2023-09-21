package keyutil

import (
	"crypto"
	"crypto/ed25519"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTripKeyGeneration(t *testing.T) {
	testDir := t.TempDir()
	// Override the source of randomness so our test is deterministic.
	randReader = rand.New(rand.NewSource(0))

	pubPath, privPath := filepath.Join(testDir, "keypair.pub"), filepath.Join(testDir, "keypair.key")
	if err := GenerateED25519ToFiles(pubPath, privPath); err != nil {
		t.Fatalf("GenerateED25519ToFiles: %v", err)
	}

	// Now, load those keys back and use them
	priv, err := DecodeED25519PrivateKeyFromFile(privPath)
	if err != nil {
		t.Fatalf("failed to decode private key: %v", err)
	}

	pub, err := DecodeED25519PublicKeyFromFile(pubPath)
	if err != nil {
		t.Fatalf("failed to decode public key: %v", err)
	}

	// Sign some data
	msg := []byte("I am a test message!")
	sig, err := priv.Sign(randReader, msg, crypto.Hash(0))
	if err != nil {
		t.Fatalf("failed to sign test message: %v", err)
	}

	// We check this to make sure our test is deterministic, and because we want to
	// make sure we're not generating all zeroes or something with impossibly bad
	// generated keys.
	want := []byte{
		0x93, 0x1d, 0x12, 0x4c, 0xfe, 0x4c, 0x17, 0xa8,
		0xd2, 0xa8, 0x67, 0xe0, 0xfb, 0xe9, 0xa8, 0x0f,
		0x65, 0x7b, 0x6a, 0x05, 0x0e, 0x94, 0xca, 0x51,
		0x7f, 0x5e, 0x8b, 0x6b, 0x06, 0x06, 0x20, 0x21,
		0xf6, 0x92, 0x59, 0xca, 0xdd, 0x8f, 0x32, 0xfa,
		0xb8, 0xfa, 0xba, 0x57, 0x57, 0x51, 0xbd, 0xf5,
		0x73, 0xca, 0x6c, 0x2b, 0x10, 0x03, 0x80, 0xc6,
		0x91, 0xb9, 0xa7, 0xfc, 0xfb, 0xf1, 0xf1, 0x0a,
	}
	if diff := cmp.Diff(want, sig); diff != "" {
		t.Errorf("unexpected signature produced (-want +got)\n%s", diff)
	}

	// Now, verify the signature.
	if !ed25519.Verify(pub, msg, sig) {
		t.Fatal("public key failed to verify signature produced by private key over mesage")
	}
}
