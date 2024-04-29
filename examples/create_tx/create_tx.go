package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/go-sdk/ec/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

// https://goplay.tools/snippet/uFtQua2P5-0
func main() {
	tx := transaction.NewTx()

	_ = tx.From(
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		0,
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		1500,
	)

	_ = tx.PayToAddress("1NRoySJ9Lvby6DuE2UQYnyT67AASwNZxGb", 1000)

	decodedWif, _ := wif.DecodeWIF("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")

	if err := tx.FillAllInputs(context.Background(), &unlocker.Getter{PrivateKey: decodedWif.PrivKey}); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("tx: %s\n", tx)
}
