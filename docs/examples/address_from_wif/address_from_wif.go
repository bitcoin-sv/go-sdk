package main

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
)

func main() {

	wif, _ := wif.DecodeWIF("Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn")

	// Print the private key
	address, _ := bscript.NewAddressFromPublicKey(wif.PrivKey.PubKey(), true)

	// Print the address, and the pubkey hash
	println(address.AddressString, address.PublicKeyHash)

}
