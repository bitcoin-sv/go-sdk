package main

import (
	"log"

	script "github.com/bsv-blockchain/go-sdk-sdk/script"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

func main() {

	priv, _ := ec.PrivateKeyFromWif("Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn")

	// Print the private key
	log.Printf("Private key: %x\n", priv.Serialize())
	address, _ := script.NewAddressFromPublicKey(priv.PubKey(), true)

	// Print the address, and the pubkey hash
	println(address.AddressString, address.PublicKeyHash)

}
