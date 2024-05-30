package main

import (
	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	script "github.com/bitcoin-sv/go-sdk/script"
)

func main() {

	wif, _ := wif.DecodeWIF("Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn")

	// Print the private key
	address, _ := script.NewAddressFromPublicKey(wif.PrivKey.PubKey(), true)

	// Print the address, and the pubkey hash
	println(address.AddressString, address.PublicKeyHash)

}
