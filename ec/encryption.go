package ec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// AES / ECIES encryption

// EncryptWithPrivateKey will encrypt the data using a given private key
func EncryptWithPrivateKey(privateKey *PrivateKey, data string) (string, error) {
	var block cipher.Block
	block, err := aes.NewCipher(privateKey.PublicKey.X.Bytes())
	if err != nil {
		return "", err
	}

	// Encrypt using bec
	encryptedData, err := Encrypt(block, []byte(data))
	if err != nil {
		return "", err
	}

	// Return the hex encoded value
	return hex.EncodeToString(encryptedData), nil
}

// DecryptWithPrivateKey is a wrapper to decrypt the previously encrypted
// information, given a corresponding private key
func DecryptWithPrivateKey(privateKey *PrivateKey, data string) (string, error) {

	// Decode the hex encoded string
	rawData, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	var block cipher.Block
	block, err = aes.NewCipher(privateKey.X.Bytes())
	if err != nil {
		return "", err
	}

	// Decrypt the data
	var decrypted []byte
	if decrypted, err = Decrypt(block, rawData); err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// EncryptShared will ECIES encrypt data and provide shared keys for decryption
func EncryptShared(user1PrivateKey *PrivateKey, user2PubKey *PublicKey, data []byte) (
	*PrivateKey, *PublicKey, []byte, error) {

	// Generate shared keys that can be decrypted by either user
	sharedPrivKey, sharedPubKey := GenerateSharedKeyPair(user1PrivateKey, user2PubKey)

	// 	var block cipher.Block
	block, err := aes.NewCipher(sharedPubKey.X.Bytes())
	if err != nil {
		return nil, nil, nil, err
	}
	// Encrypt data with shared key
	encryptedData, err := Encrypt(block, data)
	return sharedPrivKey, sharedPubKey, encryptedData, err
}

// Encrypt is an encrypt function
func Encrypt(cipherBlock cipher.Block, text []byte) ([]byte, error) {
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(cipherBlock, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt is a decrypt function
func Decrypt(cipherBlock cipher.Block, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(cipherBlock, iv)
	text := ciphertext[aes.BlockSize:]
	cfb.XORKeyStream(text, text)
	return base64.StdEncoding.DecodeString(string(text))
}

// GenerateSharedKeyPair creates shared keys that can be used to encrypt/decrypt data
// that can be decrypted by yourself (privateKey) and also the owner of the given public key
func GenerateSharedKeyPair(privateKey *PrivateKey,
	pubKey *PublicKey) (*PrivateKey, *PublicKey) {
	return PrivateKeyFromBytes(
		GenerateSharedSecret(privateKey, pubKey),
	)
}

// GenerateSharedSecret generates a shared secret based on a private key and a
// public key using Diffie-Hellman key exchange (ECDH) (RFC 4753).
// RFC5903 Section 9 states we should only return x.
func GenerateSharedSecret(privkey *PrivateKey, pubkey *PublicKey) []byte {
	x, _ := pubkey.Curve.ScalarMult(pubkey.X, pubkey.Y, privkey.D.Bytes())
	return x.Bytes()
}
