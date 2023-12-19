package pactasrv

import "testing"

func TestFormatByteSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int
		expected string
	}{
		{"TestZeroBytes", 0, ""},
		{"TestNegativeBytes", -100, ""},
		{"TestBytes", 7, "7 Bytes"},
		{"TestKilobyte", 2321, "2.32 kB"},
		{"TestMegabyte", 1004999, "1.00 MB"},
		{"TestGigabyte", 40005100000, "40.01 GB"},
		{"TestTerabyte", 100000000000000, "100.00 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatByteSize(tt.bytes); got != tt.expected {
				t.Errorf("formatByteSize(%d) = %v, want %v", tt.bytes, got, tt.expected)
			}
		})
	}
}
