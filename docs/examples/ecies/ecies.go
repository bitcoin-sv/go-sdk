package docs

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/primitives"
)

// Example of using ECIES to encrypt and decrypt data
func main() {
	// user 1
	user1Pk, _ := primitives.NewPrivateKey()

	// user 2
	user2Pk, _ := primitives.PublicKeyFromString("03121a7afe56fc8e25bca4bb2c94f35eb67ebe5b84df2e149d65b9423ee65b8b4b")

	priv, _, encryptedData, _ := primitives.EncryptShared(user1Pk, user2Pk, []byte("this is a test"))

	decryptedData, _ := primitives.DecryptWithPrivateKey(priv, hex.EncodeToString(encryptedData))

	fmt.Printf("decryptedData: %s\n", decryptedData)
}
