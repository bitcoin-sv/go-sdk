package main

import (
	"context"
	"encoding/hex"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/locker"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

const wifStr = "KycyNGiqoePhCPjGZ5m2LgpQHEMnqnd8ev4k6xueSkErtTQM2pJy"

func main() {
	key, _ := wif.DecodeWIF(wifStr)

	tx := transaction.NewTx()

	script, _ := locker.P2PKH{PubKey: key.PrivKey.PubKey()}.LockingScript()

	tx.AddOutput(&transaction.Output{
		Satoshis:      1000,
		LockingScript: script,
	})

	prevTx, _ := hex.DecodeString("11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d")
	script, _ = bscript.NewFromHex("76a9145c171f2511f5f93ac8aed7c61c676842eee4283988ac")

	p2pkhUtxo := transaction.UTXO{
		TxID:          prevTx,
		Vout:          0,
		LockingScript: script,
		Satoshis:      1500,
	}

	tx.FromUTXOs(&p2pkhUtxo)

	tx.FillInput(
		context.Background(),
		&unlocker.P2PKH{PrivateKey: key.PrivKey},
		transaction.UnlockerParams{},
	)
}
