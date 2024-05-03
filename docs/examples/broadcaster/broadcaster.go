package main

import (
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/broadcaster"
)

func main() {

	hexTx := "010000000100"
	tx, _ := transaction.NewTxFromHex(hexTx)

	success, failure := tx.Broadcast(&broadcaster.Arc{
		ApiUrl: "https://arc.gorillapool.io",
		ApiKey: "",
	})

	if failure != nil {
		panic(failure)
	}

	println(success.Message, success.Txid)
}
