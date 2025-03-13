package main

import (
	"fmt"

	ec "github.com/bsv-blockchain/go-sdk-sdk/primitives/ec"
	ecies "github.com/bsv-blockchain/go-sdk/compat/ecies"
)

// Example of using ECIES to encrypt and decrypt data for a single user
func main() {
	myPrivateKey, _ := ec.PrivateKeyFromWif("L211enC224G1kV8pyyq7bjVd9SxZebnRYEzzM3i7ZHCc1c5E7dQu")

	encryptedData, _ := ecies.EncryptSingle("hello world", myPrivateKey)

	fmt.Println(encryptedData)
	// Prints:
	// QklFMQLoYyD2A6LA9Pd342B7Z5q4agY+r674wbq6Vu2YLtVqNU5RpP1SQZNkJ22FOQt9LmXHYgMFkORAJ1nD/JVGmbmmDCx4rbYfZBVh/aa9B4imlA==

	decryptedData, _ := ecies.DecryptSingle(encryptedData, myPrivateKey)

	fmt.Printf("decryptedData: %s\n", decryptedData)
	// Prints:
	// decryptedData: hello world
}
