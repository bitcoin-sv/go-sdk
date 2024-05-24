package main

import (
	"log"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

//consider

//locking opcode:"OP_2 OP_2 OP_ADD OP_EQUAL"

//txid:="b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576"

// output_index of custom locking script in tx:0
func main() {

	unLockingOpcode := "OP_4"
	tx := transaction.NewTx()

	_ = tx.From(
		"b6832d262259ad66e1ec64dfbe85a2676089e74859d9bb217d9e0ee5d40c22a9",
		0,
		"6dc49c1a3d8241621d2150ec849833cc785e1c8cc0e152e4a97fe238d8fa6a6e",
		100,
	)

	_ = tx.AddOpReturnOutput([]byte("Custom locking script unlocked!"))

	unLockingScript, _ := bscript.NewFromASM(unLockingOpcode)

	tx.InsertInputUnlockingScript(0, unLockingScript)

	log.Println("tx: ", tx.String())
}
