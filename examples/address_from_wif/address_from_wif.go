package main

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec/wif"
)

func main() {

	wif, err := wif.DecodeWIF("Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn")
	if err != nil {
		panic(err)
	}

	// Print the private key
	address, err := bscript.NewAddressFromPublicKey(wif.PrivKey.PubKey(), true)
	if err != nil {
		panic(err)
	}

	// Print the address, and the pubkey hash
	println(address.AddressString, address.PublicKeyHash)

}
