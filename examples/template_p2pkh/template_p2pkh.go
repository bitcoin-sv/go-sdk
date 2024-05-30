package main

import (
	"encoding/hex"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template"
)

const wifStr = "KycyNGiqoePhCPjGZ5m2LgpQHEMnqnd8ev4k6xueSkErtTQM2pJy"

func main() {
	key, _ := wif.DecodeWIF(wifStr)

	tx := transaction.NewTx()

	tmpl := template.NewP2PKHTemplateFromPrivKey(key.PrivKey)
	s, _ := tmpl.Lock()

	tx.AddOutput(&transaction.TransactionOutput{
		Satoshis:      1000,
		LockingScript: s,
	})

	prevTx, _ := hex.DecodeString("11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d")
	s, _ = script.NewFromHex("76a9145c171f2511f5f93ac8aed7c61c676842eee4283988ac")

	p2pkhUtxo := transaction.UTXO{
		TxID:          prevTx,
		Vout:          0,
		LockingScript: s,
		Satoshis:      1500,
	}

	tx.FromUTXOs(&p2pkhUtxo)
	tx.Sign()
}
