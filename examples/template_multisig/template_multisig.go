package main

import (
	"context"
	"encoding/hex"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec"
	"github.com/bitcoin-sv/go-sdk/ec/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/locker"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

const wifStr = "KycyNGiqoePhCPjGZ5m2LgpQHEMnqnd8ev4k6xueSkErtTQM2pJy"
const wifStr2 = "L3xAia5eAoqiS4j4Sz2quZGVMp1ZdxfDBM8eD9frmMDZoh8WyP7m"

func main() {
	key, _ := wif.DecodeWIF(wifStr)
	key2, _ := wif.DecodeWIF(wifStr2)

	tx := transaction.NewTx()

	script, _ := locker.Multisig{PubKeys: []*ec.PublicKey{key.PrivKey.PubKey()}}.LockingScript()

	tx.AddOutput(&transaction.TransactionOutput{
		Satoshis:      1000,
		LockingScript: script,
	})

	prevTx, _ := hex.DecodeString("c358b114f15baa20bf6783714e93f5c8f036653ec50841ce6e9ee5fe2b9ddf0c")
	script, _ = bscript.NewFromHex("512102cb560e47b1ae629416b4293256443cef4427cd5e5f233a8fd2a92f1912ece4a42103da972d5d07c3abd1c30938f11e9fa78b536ec833d5ba5fc0f2146aaa7d48e48352ae")

	p2pkhUtxo := transaction.UTXO{
		TxID:          prevTx,
		Vout:          1,
		LockingScript: script,
		Satoshis:      111,
	}

	tx.FromUTXOs(&p2pkhUtxo)

	tx.FillInput(
		context.Background(),
		&unlocker.Multisig{PrivateKey: key.PrivKey},
		transaction.UnlockerParams{},
	)

	tx.FillInput(
		context.Background(),
		&unlocker.Multisig{PrivateKey: key2.PrivKey},
		transaction.UnlockerParams{},
	)
}
