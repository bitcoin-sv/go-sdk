package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

func main() {
	prevTx, _ := hex.DecodeString("11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d")
	lScript, _ := bscript.NewFromHex("76a914eb0bd5edba389198e73f8efabddfc61666969ff788ac6a0568656c6c6f")

	utxo := &transaction.UTXO{
		TxID:          prevTx,
		Vout:          0,
		LockingScript: lScript,
		Satoshis:      1500,
	}

	fmt.Printf("utxo: %+v\n", utxo)
}
