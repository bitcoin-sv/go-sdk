package main

import (
	"context"
	"log"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

// https://goplay.tools/snippet/bnsS-pA56ob
func main() {
	// Create a new transaction
	tx := transaction.NewTx()

	// Add the inputs
	err := tx.From(
		// Previous transaction ID (hex)
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		// Previous transaction output index
		0,
		// Previous transaction script (hex)
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		// Previous transaction output value in satoshis
		1500,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Add the outputs
	err = tx.PayToAddress(
		// Destination address
		"1NRoySJ9Lvby6DuE2UQYnyT67AASwNZxGb",
		// Value in satoshis
		1000,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Fill all inputs with the given private key
	decodedWif, _ := wif.DecodeWIF("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")

	// Sign the transaction
	if err := tx.FillAllInputs(
		// The default context which is empty but can be
		// used to control cancellations and timeouts.
		context.Background(),
		&unlocker.Getter{PrivateKey: decodedWif.PrivKey}); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("tx: %s\n", tx)
}
