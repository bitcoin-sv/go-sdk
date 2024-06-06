package main

import (
	"log"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template"
)

// https://goplay.tools/snippet/bnsS-pA56ob
func main() {
	priv, _ := ec.PrivateKeyFromWif("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")

	// Create a new transaction
	tx := transaction.NewTx()
	tmpl := template.NewP2PKHFromPrivKey(priv)

	// Add the inputs
	err := tx.AddInputFrom(
		// Previous transaction ID (hex)
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		// Previous transaction output index
		0,
		// Previous transaction script (hex)
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		// Previous transaction output value in satoshis
		1500,
		tmpl,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Add the outputs
	payTmpl, _ := template.NewP2PKHFromAddressString("1NRoySJ9Lvby6DuE2UQYnyT67AASwNZxGb")
	err = tx.AddOutputFromTemplate(payTmpl, 1000)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Sign the transaction
	if err := tx.Sign(); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("tx: %s\n", tx)
}
