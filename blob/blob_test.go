package blob

import "testing"

const testScheme = Scheme("test")

func TestSchemeString(t *testing.T) {
	want := "test://"
	got := testScheme.String()

	if got != want {
		t.Errorf("Scheme.String() = %q, want %q", got, want)
	}
}

func TestHasScheme(t *testing.T) {
	tests := []struct {
		desc string
		uri  string
		want bool
	}{
		{
			desc: "has scheme",
			uri:  "test://container/path/to/obj",
			want: true,
		},
		{
			desc: "does not have scheme",
			uri:  "otherscheme://container/path/to/obj",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := HasScheme(testScheme, test.uri)
			if got != test.want {
				t.Errorf("HasScheme = %t, want %t", got, test.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	got := Join(testScheme, "container", "path", "to", "obj")
	want := "test://container/path/to/obj"

	if got != want {
		t.Errorf("Join = %q, want %q", got, want)
	}
}

func TestJoin_TrailingSlash(t *testing.T) {
	got := Join(testScheme, "container", "path", "to", "obj", "" /* traiiling slash */)
	want := "test://container/path/to/obj/"

	if got != want {
		t.Errorf("Join = %q, want %q", got, want)
	}
}

func TestSplitURI(t *testing.T) {
	tests := []struct {
		desc          string
		in            string
		wantNamespace string
		wantObject    string
		wantOK        bool
	}{
		{
			desc:          "valid, has namespace and object",
			in:            "test://container/path/to/obj",
			wantNamespace: "container",
			wantObject:    "path/to/obj",
			wantOK:        true,
		},
		{
			desc:          "valid, has namespace and no object",
			in:            "test://container",
			wantNamespace: "container",
			wantObject:    "",
			wantOK:        true,
		},
		{
			desc:          "valid, has namespace (with trailing slash) and no object",
			in:            "test://container/",
			wantNamespace: "container",
			wantObject:    "",
			wantOK:        true,
		},
		{
			desc:          "invalid, wrong scheme",
			in:            "otherscheme://container/path/to/obj",
			wantNamespace: "",
			wantObject:    "",
			wantOK:        false,
		},
		{
			desc:          "invalid, no namespace",
			in:            "test://",
			wantNamespace: "",
			wantObject:    "",
			wantOK:        false,
		},
		{
			desc:          "invalid, generally malformed",
			in:            "not even a scheme!",
			wantNamespace: "",
			wantObject:    "",
			wantOK:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			gotNS, gotObj, gotOK := SplitURI(testScheme, test.in)
			if gotNS != test.wantNamespace || gotObj != test.wantObject || gotOK != test.wantOK {
				t.Errorf("SplitURI = %q, %q, %t, want %q, %q, %t", gotNS, gotObj, gotOK, test.wantNamespace, test.wantObject, test.wantOK)
			}
		})
	}
}
