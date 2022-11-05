package cmd

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		input     string
		output    string
		urlCompat bool
	}{
		{"dGhpcyBpcyBhIHRlc3Q=", "this is a test", false},
		{"ICB0aGlzIGlzIGEgc3RyaW5nIHdpdGggc3BhY2Vz", "  this is a string with spaces", false},
		{"dXJsIGNvbXBhdGlibGU=", "url compatible", true},
	}

	for _, tt := range tests {
		d := Decode(tt.input, tt.urlCompat)

		fmt.Println("WHAT", d.decoded)
		if d.decoded != tt.output {
			t.Fatalf("deocded string incorrect. got=%s",
				d.decoded)

		}

	}

}
