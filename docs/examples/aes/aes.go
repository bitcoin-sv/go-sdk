package main

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/src/src/ec"
)

func main() {
	pk, _ := ec.NewPrivateKey()

	// Encrypt using the public key of the given private key
	encryptedData, _ := ec.EncryptWithPrivateKey(pk, "this is a test")

	// Decrypt using the private key
	decryptedData, _ := ec.DecryptWithPrivateKey(pk, encryptedData)

	fmt.Printf("decryptedData: %s\n", decryptedData)
}
