package pacta

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	normalizer "github.com/dimuska139/go-email-normalizer"
)

var n = normalizer.NewNormalizer()

func CanonicalizeEmail(email string) (string, error) {
	// Performs basic validation, preventing obviously malformed email addresses, but
	// not capturing the entirety of the RFC 5322 grammar.
	split := strings.Split(email, "@")
	if len(split) != 2 {
		return "", fmt.Errorf("invalid email, wrong number of at-signs: %q", email)
	}
	if !utf8.ValidString(split[0]) {
		return "", fmt.Errorf("invalid email, non-ASCII or UTF-8 local part: %q", split[0])
	}
	if len(split[0]) == 0 || len(split[1]) == 0 {
		return "", fmt.Errorf("invalid email, empty local or domain part: %q", email)
	}
	// Performs complex, vendor-specific normalizations, preventing obvious duplicates
	// like plus-aliases, and gmail dot aliases. This will return the original if it cannot
	// parse or simplify it.
	return n.Normalize(email), nil
}

func includesSpace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
