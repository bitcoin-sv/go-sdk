package primitives

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestPrivKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{
			name: "check curve",
			key: []byte{
				0xea, 0xf0, 0x2c, 0xa3, 0x48, 0xc5, 0x24, 0xe6,
				0x39, 0x26, 0x55, 0xba, 0x4d, 0x29, 0x60, 0x3c,
				0xd1, 0xa7, 0x34, 0x7d, 0x9d, 0x65, 0xcf, 0xe9,
				0x3c, 0xe1, 0xeb, 0xff, 0xdc, 0xa2, 0x26, 0x94,
			},
		},
	}

	for _, test := range tests {
		priv, pub := PrivateKeyFromBytes(test.key)

		_, err := ParsePubKey(pub.SerialiseUncompressed())
		if err != nil {
			t.Errorf("%s privkey: %v", test.name, err)
			continue
		}

		hash := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9}
		sig, err := priv.Sign(hash)
		if err != nil {
			t.Errorf("%s could not sign: %v", test.name, err)
			continue
		}

		if !sig.Verify(hash, pub) {
			t.Errorf("%s could not verify: %v", test.name, err)
			continue
		}

		serialisedKey := priv.Serialise()
		if !bytes.Equal(serialisedKey, test.key) {
			t.Errorf("%s unexpected serialised bytes - got: %x, "+
				"want: %x", test.name, serialisedKey, test.key)
		}
	}
}

// Test vector struct
type privateTestVector struct {
	SenderPublicKey     string `json:"senderPublicKey"`
	RecipientPrivateKey string `json:"recipientPrivateKey"`
	InvoiceNumber       string `json:"invoiceNumber"`
	ExpectedPrivateKey  string `json:"privateKey"`
}

func TestBRC42PrivateVectors(t *testing.T) {
	// Determine the directory of the current test file
	_, currentFile, _, _ := runtime.Caller(0)
	testdataPath := filepath.Join(filepath.Dir(currentFile), "testdata", "BRC42.private.vectors.json")

	// Read in the file
	vectors, err := os.ReadFile(testdataPath)
	if err != nil {
		t.Fatalf("Could not read test vectors: %v", err) // use Fatalf to stop test if file cannot be read
	}
	// unmarshal the json
	var testVectors []privateTestVector
	err = json.Unmarshal(vectors, &testVectors)
	if err != nil {
		t.Errorf("Could not unmarshal test vectors: %v", err)
	}
	for i, v := range testVectors {
		t.Run("BRC42 private vector #"+strconv.Itoa(i+1), func(t *testing.T) {
			publicKey, err := PublicKeyFromString(v.SenderPublicKey)
			if err != nil {
				t.Errorf("Could not parse public key: %v", err)
			}
			privateKey, err := PrivateKeyFromString(v.RecipientPrivateKey)
			if err != nil {
				t.Errorf("Could not parse private key: %v", err)
			}
			derived, err := privateKey.DeriveChild(publicKey, v.InvoiceNumber)
			if err != nil {
				t.Errorf("Could not derive child key: %v", err)
			}

			// Convert derived private key to hex and compare
			derivedHex := hex.EncodeToString(derived.Serialise())
			if derivedHex != v.ExpectedPrivateKey {
				t.Errorf("Derived private key does not match expected: got %v, want %v", derivedHex, v.ExpectedPrivateKey)
			}
		})
	}
}
