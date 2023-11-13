// Package blob provides basic primitives and helpers for working with blob
// storage systems.
package blob

import "strings"

type Scheme string

func (s Scheme) String() string {
	return string(s) + "://"
}

func Join(s Scheme, parts ...string) string {
	return s.String() + strings.Join(parts, "/")
}

func HasScheme(s Scheme, uri string) bool {
	return strings.HasPrefix(uri, s.String())
}

func SplitURI(s Scheme, uri string) (string, string, bool) {
	if !HasScheme(s, uri) {
		return "", "", false
	}
	ns, obj, ok := strings.Cut(strings.TrimPrefix(uri, s.String()), "/")
	if !ok {
		return ns, "", ns != ""
	}
	return ns, obj, true
}
