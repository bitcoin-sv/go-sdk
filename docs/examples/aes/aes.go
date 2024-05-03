package docs

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/primitives"
)

func main() {
	pk, _ := primitives.NewPrivateKey()

	// Encrypt using the public key of the given private key
	encryptedData, _ := primitives.EncryptWithPrivateKey(pk, "this is a test")

	// Decrypt using the private key
	decryptedData, _ := primitives.DecryptWithPrivateKey(pk, encryptedData)

	fmt.Printf("decryptedData: %s\n", decryptedData)
}
