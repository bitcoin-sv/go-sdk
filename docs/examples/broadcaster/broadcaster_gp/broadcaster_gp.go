package main

import (
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/broadcaster"
)

func main() {

	// Create a new transaction
	hexTx := "010000000100"
	tx, _ := transaction.NewTxFromHex(hexTx)

	// Broadcast the transaction
	success, failure := tx.Broadcast(&broadcaster.Arc{
		ApiUrl: "https://arc.gorillapool.io",
		ApiKey: "",
	})

	// Check for errors
	if failure != nil {
		panic(failure)
	}

	// Print the success message and transaction ID
	println(success.Message, success.Txid)
}
