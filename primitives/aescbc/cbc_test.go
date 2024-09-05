package primitives

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAESCBCEncryptDecrypt(t *testing.T) {
	validTestData := "Test data"
	// Test cases
	testCases := []struct {
		name        string
		data        string
		key         string
		iv          string
		expectError bool
	}{
		{
			name:        "Valid short message",
			data:        validTestData,
			key:         "0123456789abcdef0123456789abcdef", // 32 bytes (256 bits)
			iv:          "0123456789abcdef0123456789abcdef", // 32 hex chars = 16 bytes
			expectError: false,
		},
		{
			name:        "Valid long message",
			data:        "This is a longer message that spans multiple AES blocks.",
			key:         "0123456789abcdef0123456789abcdef", // 32 bytes (256 bits)
			iv:          "fedcba9876543210fedcba9876543210", // 32 hex chars = 16 bytes
			expectError: false,
		},
		{
			name:        "Invalid key length",
			data:        validTestData,
			key:         "0123456789abcdef", // 16 bytes (128 bits) - too short for AES-256
			iv:          "0123456789abcdef0123456789abcdef",
			expectError: true,
		},
		{
			name:        "Invalid IV length",
			data:        validTestData,
			key:         "0123456789abcdef0123456789abcdef",
			iv:          "211234560123456789abcdefababababab", // More than 16 bytes
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, _ := hex.DecodeString(tc.key)

			iv, _ := hex.DecodeString(tc.iv)

			// Check IV length
			if len(iv) != aes.BlockSize {
				if tc.expectError {
					return // Expected error due to invalid IV length
				}
				t.Fatalf("IV length must be %d bytes, got %d bytes", aes.BlockSize, len(iv))
			}

			data := []byte(tc.data)

			// Encrypt
			encrypted, err := AESCBCEncrypt(data, key, iv, false)
			if tc.expectError {
				require.Error(t, err, "Expected an error but got none")
				return
			}
			require.NoError(t, err, "Encryption failed")

			// Decrypt
			decrypted, err := AESCBCDecrypt(encrypted, key, iv)
			require.NoError(t, err, "Decryption failed")

			// Compare
			require.Equal(t, data, decrypted, "Decrypted data doesn't match original")
		})
	}
}

func TestPKCS7Padding(t *testing.T) {
	testCases := []struct {
		name      string
		data      string
		blockSize int
	}{
		{"Short data", "Hello", 16},
		{"Block-sized data", "1234567890123456", 16},
		{"Long data", "This is a longer string for testing", 16},
		{"Empty data", "", 16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			padded := PKCS7Padd(data, tc.blockSize)

			// Check if padded data length is multiple of block size
			require.Equal(t, 0, len(padded)%tc.blockSize, "Padded data length is not a multiple of block size")

			// Strip padding
			stripped, err := PKCS7Unpad(padded, tc.blockSize)
			require.NoError(t, err, "StripPKCS7Padding failed")

			// Compare stripped data with original
			require.Equal(t, data, stripped, "Stripped data doesn't match original")
		})
	}
}

func TestStripPKCS7PaddingInvalidPadding(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		blockSize int
	}{
		{"Invalid length", []byte{1, 2, 3}, 16},
		{"Large padding byte", []byte{1, 2, 3, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17}, 16},
		{"Inconsistent padding", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 4}, 16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := PKCS7Unpad(tc.data, tc.blockSize)
			require.Error(t, err, "Expected an error, but got nil")
		})
	}
}

func TestAESEncryptDecryptWithRandomIV(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err, "Failed to generate random key")

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	require.NoError(t, err, "Failed to generate random IV")

	data := []byte("This is a test message with random key and IV.")

	encrypted, err := AESCBCEncrypt(data, key, iv, false)
	require.NoError(t, err, "Encryption failed")

	decrypted, err := AESCBCDecrypt(encrypted, key, iv)
	require.NoError(t, err, "Decryption failed")

	require.Equal(t, data, decrypted, "Decrypted data doesn't match original")
}
