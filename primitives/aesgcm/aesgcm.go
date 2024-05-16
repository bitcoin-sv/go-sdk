package primitives

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// AESEncrypt performs AES block encryption without any mode of operation
func AESEncrypt(plaintext, key []byte) ([]byte, error) {
	if len(plaintext) != aes.BlockSize {
		return nil, fmt.Errorf("plaintext is not the correct block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize)
	block.Encrypt(ciphertext, plaintext)

	return ciphertext, nil
}

// AESDecrypt performs AES block decryption without any mode of operation
func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) != aes.BlockSize {
		return nil, fmt.Errorf("ciphertext is not the correct block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, aes.BlockSize)
	block.Decrypt(plaintext, ciphertext)

	return plaintext, nil
}

// EncryptGCM encrypts plaintext using AES-GCM with the provided key and additional data
func EncryptGCM(plaintext, nonce, key, additionalData []byte) (ciphertext, tag []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Automatically adjust GCM to the nonce size
	gcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, additionalData)
	tag = ciphertext[len(ciphertext)-gcm.Overhead():]

	return ciphertext[:len(ciphertext)-gcm.Overhead()], tag, nil
}

// DecryptGCM decrypts ciphertext using AES-GCM with the provided key, nonce, additional data, and tag
func DecryptGCM(ciphertext, key, nonce, additionalData, tag []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Automatically adjust GCM to the nonce size
	gcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return nil, err
	}

	ciphertextWithTag := append(ciphertext, tag...)
	plaintext, err = gcm.Open(nil, nonce, ciphertextWithTag, additionalData)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return plaintext, nil
}
