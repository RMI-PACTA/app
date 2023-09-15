package pacta

import (
	"fmt"
	"testing"
)

func TestEmailCanonicalization(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		// Testing valid Gmail addresses with different variations
		{"john.doe@gmail.com", "johndoe@gmail.com", false},
		{"john.doe+123@gmail.com", "johndoe@gmail.com", false},
		{"j.o.h.n.d.o.e@gmail.com", "johndoe@gmail.com", false},
		{"johndoe+test@gmail.com", "johndoe@gmail.com", false},
		{"john.doe@googlemail.com", "johndoe@gmail.com", false},

		// Testing valid non-Gmail addresses
		{"johndoe@protonmail.ch", "johndoe@protonmail.ch", false},
		{"john+doe@protonmail.ch", "johndoe@protonmail.ch", false},
		{"john.doe+test@outlook.com", "john.doetest@outlook.com", false},
		{"john@xn--938h.com", "john@xn--938h.com", false},

		// Testing spacing + unicode
		{"john doe@gmail.com", "john doe@gmail.com", false},
		{"john🙂doe@gmail.com", "john🙂doe@gmail.com", false},

		// Testing invalid email addresses
		{in: "john.doe", wantErr: true},
		{in: "john.doe@", wantErr: true},
		{in: "@gmail.com", wantErr: true},
		{in: "@gmail.com", wantErr: true},
		{in: "johndoe@gmail	.com", wantErr: true},
		{in: "john@🙂.com", wantErr: true},

		// Testing valid email addresses with different domains
		{"john.doe@hotmail.com", "john.doe@hotmail.com", false},
		{"j.o.h.n.doe+test@domain.com", "j.o.h.n.doe+test@domain.com", false},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			got, err := CanonicalizeEmail(test.in)
			if test.wantErr {
				if err == nil {
					t.Fatalf("CanonicalizeEmail(%q) returned no error, but one was expected", test.in)
				}
				return
			}
			if err != nil {
				t.Errorf("CanonicalizeEmail: %v", err)
			}
			if got != test.want {
				t.Errorf("CanonicalizeEmail(%q) = %q, want %q", test.in, got, test.want)
			}
		})
	}
}
