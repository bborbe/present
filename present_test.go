package main

import (
	"testing"
	"bytes"
)

func TestRead(t *testing.T) {
	var tests = []struct {
		name           string
		content        string
		expectedErr    error
		expectedLength int
	}{
		{"empty", "", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.content)
			results, err := read(r)
			if err != tt.expectedErr {
				t.Fatalf("expected error %v got %v", tt.expectedErr, err)
			}
			if len(results) != tt.expectedLength {
				t.Fatalf("expected length %d got %d", tt.expectedLength, len(results))
			}
		})
	}
}
