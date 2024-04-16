package main

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/transaction"
)

func main() {
	utxo, _ := transaction.NewUTXO(
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		0,
		"76a914eb0bd5edba389198e73f8efabddfc61666969ff788ac6a0568656c6c6f",
		1500,
	)

	fmt.Printf("utxo: %+v\n", utxo)
}
