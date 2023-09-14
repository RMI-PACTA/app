package pacta

import (
	"fmt"
	"testing"
)

func TestEmailCanonization(t *testing.T) {
	cases := []struct {
		input       string
		ouptut      string
		errExpected bool
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
		{"john@🙂.com", "john@🙂.com", false},

		// Testing invalid email addresses
		{"john.doe", "", true},
		{"john.doe@", "", true},
		{"john doe@gmail.com", "", true},
		{"johndoe@gmail	.com", "", true},
		{"@gmail.com", "", true},
		{"@gmail.com", "", true},
		{"john🙂doe@gmail.com", "", true},

		// Testing valid email addresses with different domains
		{"john.doe@hotmail.com", "john.doe@hotmail.com", false},
		{"j.o.h.n.doe+test@domain.com", "j.o.h.n.doe+test@domain.com", false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			output, err := CanonizeEmail(c.input)
			if c.errExpected && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !c.errExpected && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if output != c.ouptut {
				t.Errorf("expected %q, got %q", c.ouptut, output)
			}
		})
	}
}
