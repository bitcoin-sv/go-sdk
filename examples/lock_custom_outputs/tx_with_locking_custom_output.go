package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

func main() {
	lockedAmount := 10
	lockingOpcode := "OP_2 OP_2 OP_ADD OP_EQUAL"
	tx := transaction.NewTx()

	_ = tx.From(
		"4cb9087f9e773c4f3a00e1209fc09c9d2fb656f3c0160dd991a55f1bc31a0bba",
		0,
		"76a9146d44da8f25965f7b5d298e19cf95dbfb9fd15ab888ac",
		200,
	)

	lockingScriptFromOpcode, err := bscript.NewFromASM(lockingOpcode)
	if err != nil {
		return
	}

	output := transaction.Output{
		Satoshis:      uint64(lockedAmount),
		LockingScript: lockingScriptFromOpcode,
	}

	tx.AddOutput(&output)

	tx.ChangeToAddress("1Axm8xCYyEDmuWTdWWaG1CkPuk8DV6xHVE", transaction.NewFeeQuote())

	decodedWif, _ := wif.DecodeWIF("L5NGkLTDmzS6KfaBHj7CUsx6FwCysMoqiYyG8fVYM6e44rMFg8DB")

	err = tx.FillAllInputs(context.Background(), &unlocker.Getter{PrivateKey: decodedWif.PrivKey})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("tx: ", tx.String())
}
