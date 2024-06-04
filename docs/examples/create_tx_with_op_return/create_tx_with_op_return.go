package main

import (
	"encoding/hex"
	"log"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template"
)

func main() {
	w, _ := wif.DecodeWIF("L3VJH2hcRGYYG6YrbWGmsxQC1zyYixA82YjgEyrEUWDs4ALgk8Vu")

	tx := transaction.NewTx()
	tmpl := template.NewP2PKHFromPrivKey(w.PrivKey)

	txid, _ := hex.DecodeString("b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576")
	s, _ := script.NewFromHex("76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac")
	tx.AddInputWithOutput(&transaction.TransactionInput{
		SourceTXID:       txid,
		SourceTxOutIndex: 0,
		Template:         tmpl,
	}, &transaction.TransactionOutput{
		LockingScript: s,
		Satoshis:      1000,
	})
	_ = tx.AddOpReturnOutput([]byte("You are using go-sdk!"))

	if err := tx.Sign(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("tx: ", tx.String())
}
