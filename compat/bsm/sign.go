package bsm

import (
	"bytes"
	"encoding/base64"
	"errors"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

const hBSV = "Bitcoin Signed Message:\n"

// SignMessage signs a string with the provided PrivateKey using Bitcoin Signed Message encoding
// sigRefCompressedKey bool determines whether the signature will reference a compressed or uncompresed key
// Spec: https://github.com/bitcoin/bitcoin/pull/524
func SignMessage(privateKey *ec.PrivateKey, message string, sigRefCompressedKey bool) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is required")
	}

	if len(message) == 0 {
		return "", errors.New("message is required")
	}

	b := new(bytes.Buffer)
	var err error

	varInt := transaction.VarInt(len(hBSV))
	b.Write(varInt.Bytes())

	// append the hBsv to buff
	b.WriteString(hBSV)

	varInt = transaction.VarInt(len(message))
	b.Write(varInt.Bytes())

	// append the data to buff
	b.WriteString(message)

	// Create the hash
	messageHash := crypto.Sha256d(b.Bytes())

	// Sign
	var sigBytes []byte
	if sigBytes, err = ec.SignCompact(ec.S256(), privateKey, messageHash, sigRefCompressedKey); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sigBytes), nil
}
