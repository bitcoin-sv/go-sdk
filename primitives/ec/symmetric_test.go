package primitives

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSymmetricKeyEncryptionAndDecryption(t *testing.T) {
	t.Logf("Running encryption and decryption without errors")
	symmetricKey := NewSymmetricKeyFromRandom()
	cipherText, err := symmetricKey.Encrypt([]byte("a thing to encrypt"))
	if err != nil {
		t.Errorf("Error encrypting: %v", err)
	}

	decrypted, err := symmetricKey.Decrypt(cipherText)
	if err != nil {
		t.Errorf("Error decrypting: %v", err)
	}

	if string(decrypted) != "a thing to encrypt" {
		t.Errorf("Decrypted value does not match original plaintext")
	}
}

type symmetricTestVector struct {
	Key        string `json:"key"`
	Plaintext  string `json:"plaintext"`
	Ciphertext string `json:"ciphertext"`
}

func TestSymmetricKeyDecryption(t *testing.T) {
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
		t.Logf("Running decryption test vector %d", i+1)

		vectorCiphertext, err := base64.StdEncoding.DecodeString(v.Ciphertext)
		if err != nil {
			log.Fatalf("Failed to decode ciphertext: %v", err)
		}

		symmetricKey := NewSymmetricKeyFromString(v.Key)
		decrypted, err := symmetricKey.Decrypt(vectorCiphertext)
		if err != nil {
			t.Errorf("Error decrypting: %v", err)
		}

		if string(decrypted) != v.Plaintext {
			t.Errorf("Decrypted value does not match expected plaintext")
		}

	}
}
