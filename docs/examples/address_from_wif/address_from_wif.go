package docs

import (
	"github.com/jonesjBSV/go-sdk/script"
	"github.com/jonesjBSV/go-sdk/primitives/wif"
)

func main() {

	wif, _ := wif.DecodeWIF("Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn")

	// Print the private key
	address, _ := bscript.NewAddressFromPublicKey(wif.PrivKey.PubKey(), true)

	// Print the address, and the pubkey hash
	println(address.AddressString, address.PublicKeyHash)

}
