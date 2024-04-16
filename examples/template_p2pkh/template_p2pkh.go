package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec/wif"
	"github.com/bitcoin-sv/go-sdk/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/txbuilder"
)

const wifStr = "KycyNGiqoePhCPjGZ5m2LgpQHEMnqnd8ev4k6xueSkErtTQM2pJy"

func main() {
	key, _ := wif.DecodeWIF(wifStr)

	tx := transaction.NewTx()

	out := &txbuilder.P2PKH{
		PubKey: key.PrivKey.PubKey(),
	}

	scr, _ := out.LockingScript()

	tx.AddOutput(&transaction.Output{
		Satoshis:      1000,
		LockingScript: scr,
	})

	prevTx, _ := hex.DecodeString("11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d")
	script, _ := bscript.NewFromHex("76a914eb0bd5edba389198e73f8efabddfc61666969ff788ac6a0568656c6c6f")
	scriptP2pkh := &txbuilder.P2PKH{
		Script: *script,
	}

	fmt.Println(scriptP2pkh.IsLockingScript())

	p2pkhUtxo := transaction.UTXO{
		TxID:          prevTx,
		Vout:          0,
		LockingScript: script,
		Satoshis:      1500,
	}

	tx.FromUTXOs(&p2pkhUtxo)

	unlock, _ := scriptP2pkh.UnlockingScript(*tx, transaction.UnlockerParams{
		InputIdx:     0,
		SigHashFlags: sighash.AllForkID,
	})

	tx.Inputs[0].UnlockingScript = unlock

}
