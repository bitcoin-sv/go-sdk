package main

import (
	"fmt"

	message "github.com/bitcoin-sv/go-sdk/message"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

func main() {

	fmt.Printf("encrypt / decrypt using keys")

	// sender private key
	pk, err := ec.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	// receiver private key
	pkr, err := ec.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	// Encrypt using the public key of the given private key
	encryptedData, err := message.Encrypt([]byte("this is a test"), pk, pkr.PubKey())

	// Decrypt using the private key
	decryptedData, err := message.Decrypt(encryptedData, pkr)

	fmt.Printf("decryptedData: %s\n", decryptedData)

}
