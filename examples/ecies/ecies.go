package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/ec"
)

// Example of using ECIES to encrypt and decrypt data
func main() {
	// user 1
	user1Pk, err := ec.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	// user 2
	user2Pk, err := ec.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")
	if err != nil {
		panic(err)
	}

	priv, _, encryptedData, err := ec.EncryptShared(user1Pk, user2Pk, []byte("this is a test"))
	if err != nil {
		panic(err)
	}

	decryptedData, err := ec.DecryptWithPrivateKey(priv, hex.EncodeToString(encryptedData))
	if err != nil {
		panic(err)
	}

	fmt.Printf("decryptedData: %s\n", decryptedData)

}
