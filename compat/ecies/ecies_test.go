package compat

import (
	"encoding/base64"
	"log"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/stretchr/testify/require"
)

const msgString = "hello world"

func TestEncryptDecryptSingle(t *testing.T) {

	pk, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	// Electrum Encrypt
	encryptedData, err := EncryptSingle(msgString, pk)
	require.NoError(t, err)
	require.Equal(t, "QklFMQO7zpX/GS4XpthCy6/hT38ZKsBGbn8JKMGHOY5ifmaoT890Krt9cIRk/ULXaB5uC08owRICzenFbm31pZGu0gCM2uOxpofwHacKidwZ0Q7aEw==", encryptedData)

	// Electrum Decrypt
	decryptedData, _ := DecryptSingle(encryptedData, pk)
	require.Equal(t, msgString, decryptedData)
}

func TestElectrumEncryptDecryptSingle(t *testing.T) {

	pk, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	// Electrum Encrypt
	encryptedData, err := ElectrumEncrypt([]byte(msgString), pk.PubKey(), pk, false)
	require.NoError(t, err)
	expectedB64, err := base64.StdEncoding.DecodeString("QklFMQO7zpX/GS4XpthCy6/hT38ZKsBGbn8JKMGHOY5ifmaoT890Krt9cIRk/ULXaB5uC08owRICzenFbm31pZGu0gCM2uOxpofwHacKidwZ0Q7aEw==")
	require.NoError(t, err)
	require.Equal(t, expectedB64, encryptedData)

	// Electrum Decrypt
	decryptedData, _ := ElectrumDecrypt(encryptedData, pk, nil)
	require.Equal(t, []byte(msgString), decryptedData)
}

func TestElectrumEncryptDecryptSingleNokey(t *testing.T) {

	pk, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	// Electrum Encrypt
	encryptedData, err := ElectrumEncrypt([]byte(msgString), pk.PubKey(), pk, true)
	require.NoError(t, err)

	// Electrum Decrypt
	decryptedData, _ := ElectrumDecrypt(encryptedData, pk, pk.PubKey())
	require.Equal(t, []byte(msgString), decryptedData)
}

func TestElectrumEncryptDecryptWithCounterparty(t *testing.T) {
	pk1, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")
	counterparty, err := ec.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")
	require.NoError(t, err)

	// Electrum Encrypt
	encryptedData, _ := ElectrumEncrypt([]byte(msgString), counterparty, pk1, false)
	require.NoError(t, err)
	log.Println("Encrypted data: ", encryptedData)

	// Electrum Decrypt
	decryptedData, err := ElectrumDecrypt(encryptedData, pk1, counterparty)
	require.NoError(t, err)
	require.Equal(t, msgString, string(decryptedData))
}

func TestElectrumEncryptDecryptWithCounterpartyNoKey(t *testing.T) {
	pk1, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")
	counterparty, err := ec.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")
	require.NoError(t, err)

	// Electrum Encrypt
	encryptedData, _ := ElectrumEncrypt([]byte(msgString), counterparty, pk1, true)
	require.NoError(t, err)
	log.Println("Encrypted data: ", encryptedData)

	// Electrum Decrypt
	decryptedData, err := ElectrumDecrypt(encryptedData, pk1, counterparty)
	require.NoError(t, err)
	require.Equal(t, msgString, string(decryptedData))
}

func TestBitcoreEncryptDecryptSolo(t *testing.T) {
	expectedEncryptedData := "A7vOlf8ZLhem2ELLr+FPfxkqwEZufwkowYc5jmJ+ZqhPAAAAAAAAAAAAAAAAAAAAAB27kUY/HpNbiwhYSpEoEZZDW+wEjMmPNcAAxnc0kiuQ73FpFzf6p6afe4wwVtKAAg=="
	decodedExpectedEncryptedData, _ := base64.StdEncoding.DecodeString(expectedEncryptedData)
	log.Printf("Decoded expected encrypted data: %x\n", decodedExpectedEncryptedData)
	pk, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	// Bitcore Encrypt
	encryptedData, err := BitcoreEncrypt([]byte(msgString), pk.PubKey(), pk, nil)
	require.NoError(t, err)
	require.Equal(t, decodedExpectedEncryptedData, encryptedData)

	// Bitcore Decrypt
	decryptedData, err := BitcoreDecrypt(encryptedData, pk)
	require.NoError(t, err)
	require.Equal(t, msgString, string(decryptedData))
}

func TestBitcoreEncryptDecryptWithCounterparty(t *testing.T) {
	pk1, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")
	counterpartyPk, err := ec.PrivateKeyFromWif("L27ZSAC1xTsZrghYHqnxwAQZ12bH57piaAdoGaLizTp3JZrjkZjK")
	require.NoError(t, err)

	// Bitcore Encrypt
	encryptedData, err := BitcoreEncrypt([]byte(msgString), counterpartyPk.PubKey(), pk1, nil)
	expectedEncryptedData := "A7vOlf8ZLhem2ELLr+FPfxkqwEZufwkowYc5jmJ+ZqhPAAAAAAAAAAAAAAAAAAAAAAmFslNpNc4TrjaMPmPLdooZwoP6/fE7GN3AeyLpFf2f+QGYRKIke8zbhxu8FcLOsA=="
	decodedExpectedEncryptedData, _ := base64.StdEncoding.DecodeString(expectedEncryptedData)

	require.NoError(t, err)
	require.Equal(t, decodedExpectedEncryptedData, encryptedData)

	// Bitcore Decrypt
	decryptedData, err := BitcoreDecrypt(encryptedData, counterpartyPk)
	require.NoError(t, err)
	require.Equal(t, msgString, string(decryptedData))
}
