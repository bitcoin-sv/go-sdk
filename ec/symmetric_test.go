package ec_test

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bitcoin-sv/go-sdk/ec"
)

type symmetricTestVector struct {
	Key        string `json:"key"`
	Plaintext  string `json:"plaintext"`
	Ciphertext string `json:"ciphertext"`
}

func TestSymmetricKeyEncryptionAndDecryption(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	testdataPath := filepath.Join(filepath.Dir(currentFile), "testdata", "SymmetricKey.vectors.json")

	vectors, err := os.ReadFile(testdataPath)
	if err != nil {
		t.Fatalf("Error reading test vectors: %v", err)
	}

	var testVectors []symmetricTestVector

	err = json.Unmarshal(vectors, &testVectors)
	if err != nil {
		t.Fatalf("Error unmarshalling test vectors: %v", err)
	}

	for i, v := range testVectors {
		expectedPlaintext, err := base64.StdEncoding.DecodeString(v.Plaintext)
		if err != nil {
			log.Fatalf("Failed to decode plaintext: %v", err)
		}
		expectedCiphertext, err := base64.StdEncoding.DecodeString(v.Ciphertext)
		if err != nil {
			log.Fatalf("Failed to decode ciphertext: %v", err)
		}

		t.Run("Encryption", func(t *testing.T) {
			t.Logf("Running encryption test vector %d", i+1)
			symmetricKey := ec.NewSymmetricKeyFromString(v.Key)
			cipherText, err := symmetricKey.Encrypt(expectedPlaintext)
			if err != nil {
				t.Errorf("Error encrypting: %v", err)
			}

			if string(cipherText) != v.Ciphertext {
				t.Errorf("Encrypted value does not match expected ciphertext")
			}
		})

		t.Run("Decryption", func(t *testing.T) {
			t.Logf("Running decryption test vector %d", i+1)
			symmetricKey := ec.NewSymmetricKeyFromString(v.Key)
			decrypted, err := symmetricKey.Decrypt(expectedCiphertext)
			if err != nil {
				t.Errorf("Error decrypting: %v", err)
			}

			if string(decrypted) != v.Plaintext {
				t.Errorf("Decrypted value does not match expected plaintext")
			}
		})

		t.Run("End to end", func(t *testing.T) {
			t.Logf("Running encryption and decryption without errors %d", i+1)
			symmetricKey := ec.NewSymmetricKeyFromString(v.Key)
			cipherText, err := symmetricKey.Encrypt(expectedPlaintext)
			if err != nil {
				t.Errorf("Error encrypting: %v", err)
			}

			_, err = symmetricKey.Decrypt(cipherText)
			if err != nil {
				t.Errorf("Error decrypting: %v", err)
			}
		})
	}
}
