package compat

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBase58(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		encoded string
		err     bool
	}{
		{
			name:    "empty string",
			input:   "",
			encoded: "",
			err:     false,
		},
		{
			name:    "single zero byte",
			input:   "00",
			encoded: "1",
			err:     false,
		},
		{
			name:    "decoded address",
			input:   "00010966776006953D5567439E5E39F86A0D273BEED61967F6",
			encoded: "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
			err:     false,
		},
		{
			name:    "decoded hash",
			input:   "0123456789ABCDEF",
			encoded: "C3CPq7c8PY",
			err:     false,
		},
		{
			name:    "leading zeros",
			input:   "000000287FB4CD",
			encoded: "111233QC4",
			err:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			inputBytes, err := hex.DecodeString(tt.input)
			if err != nil {
				t.Fatalf("failed to decode hex input: %v", err)
			}

			encoded := Encode(inputBytes)
			if encoded != tt.encoded {
				t.Errorf("Encode() = %v, want %v", encoded, tt.encoded)
			}

			// Test decoding
			decoded, err := Decode(tt.encoded)
			if (err != nil) != tt.err {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !tt.err && !bytes.Equal(decoded, inputBytes) {
				t.Errorf("Decode() = %x, want %x", decoded, inputBytes)
			}
		})
	}
}

func TestBase58DecodeInvalid(t *testing.T) {
	invalidInputs := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid character",
			input: "invalid!@#$%",
		},
		{
			name:  "mixed valid and invalid",
			input: "1234!@#$%",
		},
	}

	for _, tt := range invalidInputs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			if err == nil {
				t.Error("Decode() expected error for invalid input")
			}
		})
	}
}

func TestBase58EncodeEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "nil input",
			input: nil,
			want:  "",
		},
		{
			name:  "all zeros",
			input: []byte{0, 0, 0, 0},
			want:  "1111",
		},
		{
			name:  "large number",
			input: []byte{255, 255, 255, 255},
			want:  "7YXq9G",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
