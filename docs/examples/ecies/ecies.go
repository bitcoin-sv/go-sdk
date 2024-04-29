package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/src/ec"
)

// Example of using ECIES to encrypt and decrypt data
func main() {
	// user 1
	user1Pk, _ := ec.NewPrivateKey()

	// user 2
	user2Pk, _ := ec.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")

	priv, _, encryptedData, _ := ec.EncryptShared(user1Pk, user2Pk, []byte("this is a test"))

	decryptedData, _ := ec.DecryptWithPrivateKey(priv, hex.EncodeToString(encryptedData))

	fmt.Printf("decryptedData: %s\n", decryptedData)
}
