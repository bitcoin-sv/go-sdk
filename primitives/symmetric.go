package primitives

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/bitcoin-sv/go-sdk/aesgcm"
)

type SymmetricKey struct {
	key []byte
}

// Encrypt encrypts the given message using the symmetric key using AES-GCM
func (s *SymmetricKey) Encrypt(message []byte) (ciphertext []byte, err error) {
	iv := make([]byte, 32)
	rand.Read(iv)
	cipertext, tag, err := aesgcm.EncryptGCM(message, iv, s.ToBytes(), []byte{})
	if err != nil {
		return nil, err
	}
	return append(append(iv, cipertext...), tag...), nil
}

// Decrypt decrypts the given message using the symmetric key using AES-GCM
func (s *SymmetricKey) Decrypt(message []byte) (plaintext []byte, err error) {
	iv := message[:32]
	ciphertext := message[32 : len(message)-16]
	tag := message[len(message)-16:]
	plaintext, err = aesgcm.DecryptGCM(ciphertext, s.ToBytes(), iv, []byte{}, tag)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (s *SymmetricKey) ToBytes() []byte {
	return s.key
}

func (s *SymmetricKey) FromBytes(b []byte) *SymmetricKey {
	return &SymmetricKey{key: b}
}

func NewSymmetricKey(key []byte) *SymmetricKey {
	return &SymmetricKey{key: key}
}

func NewSymmetricKeyFromRandom() *SymmetricKey {
	key := make([]byte, 32)
	rand.Read(key)
	return &SymmetricKey{key: key}
}

func NewSymmetricKeyFromString(keyBase64String string) *SymmetricKey {
	// Decode the Base64 string to bytes
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64String)
	if err != nil {
		log.Fatalf("Failed to decode Base64 symmetric key string: %v", err)
	}
	return &SymmetricKey{key: keyBytes}
}
