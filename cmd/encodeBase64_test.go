package cmd

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		input     string
		output    string
		urlCompat bool
	}{
		{"this is a test", "dGhpcyBpcyBhIHRlc3Q=", false},
		{"  this is a string with spaces", "ICB0aGlzIGlzIGEgc3RyaW5nIHdpdGggc3BhY2Vz", false},
		{"url compatible", "dXJsIGNvbXBhdGlibGU=", true},
	}

	for _, tt := range tests {
		d := Encode(tt.input, tt.urlCompat)

		if d.encoded != tt.output {
			t.Fatalf("enocded string incorrect. got=%s",
				d.encoded)

		}

	}

}
