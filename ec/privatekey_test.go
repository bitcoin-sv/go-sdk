package ec

import (
	"bytes"
	"encoding/hex"
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
		priv, pub := PrivateKeyFromBytes(S256(), test.key)

		_, err := ParsePubKey(pub.SerialiseUncompressed(), S256())
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
type testVector struct {
	senderPublicKey     string
	recipientPrivateKey string
	invoiceNumber       string
	expectedPrivateKey  string
}

// Test vectors
var testVectors = []testVector{
	{
		senderPublicKey:     "033f9160df035156f1c48e75eae99914fa1a1546bec19781e8eddb900200bff9d1",
		recipientPrivateKey: "6a1751169c111b4667a6539ee1be6b7cd9f6e9c8fe011a5f2fe31e03a15e0ede",
		invoiceNumber:       "f3WCaUmnN9U=",
		expectedPrivateKey:  "761656715bbfa172f8f9f58f5af95d9d0dfd69014cfdcacc9a245a10ff8893ef",
	},
	// Add the remaining vectors...
}

func TestBRC42Vectors(t *testing.T) {
	for i, v := range testVectors {
		t.Run("BRC42 private vector #"+strconv.Itoa(i+1), func(t *testing.T) {
			publicKey := PublicKeyFromString(v.senderPublicKey)
			privateKey, err := PrivateKeyFromString(v.recipientPrivateKey)
			if err != nil {
				t.Errorf("Could not parse private key: %v", err)
			}
			derived := privateKey.DeriveChild(publicKey, v.invoiceNumber)

			// Convert derived private key to hex and compare
			derivedHex := hex.EncodeToString(derived.Serialise())
			if derivedHex != v.expectedPrivateKey {
				t.Errorf("Derived private key does not match expected: got %v, want %v", derivedHex, v.expectedPrivateKey)
			}
		})
	}
}
