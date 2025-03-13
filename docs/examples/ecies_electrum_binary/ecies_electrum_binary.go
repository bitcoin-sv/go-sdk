package main

import (
	"fmt"

	ec "github.com/bsv-blockchain/go-sdk-sdk/primitives/ec"
	ecies "github.com/bsv-blockchain/go-sdk/compat/ecies"
)

// Example of using ECIES to encrypt and decrypt data
func main() {

	// user 1
	user1Pk, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	// user 2
	user2Pk, _ := ec.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")

	encryptedData, _ := ecies.ElectrumEncrypt([]byte("hello world"), user2Pk, user1Pk, false)

	fmt.Println(encryptedData)
	decryptedData, _ := ecies.ElectrumDecrypt(encryptedData, user1Pk, user2Pk)

	fmt.Printf("decryptedData: %s\n", decryptedData)

}
