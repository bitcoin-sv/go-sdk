package main

import (
	"encoding/hex"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template"
)

const wifStr = "KycyNGiqoePhCPjGZ5m2LgpQHEMnqnd8ev4k6xueSkErtTQM2pJy"
const wifStr2 = "L3xAia5eAoqiS4j4Sz2quZGVMp1ZdxfDBM8eD9frmMDZoh8WyP7m"

func main() {
	key, _ := wif.DecodeWIF(wifStr)
	key2, _ := wif.DecodeWIF(wifStr2)

	tx := transaction.NewTx()

	tmpl := template.NewMultisigTemplateFromPrivKeys([]*ec.PrivateKey{key.PrivKey, key2.PrivKey}, 1)
	s, _ := tmpl.Lock()

	tx.AddOutput(&transaction.TransactionOutput{
		Satoshis:      1000,
		LockingScript: s,
	})

	prevTx, _ := hex.DecodeString("c358b114f15baa20bf6783714e93f5c8f036653ec50841ce6e9ee5fe2b9ddf0c")
	s, _ = script.NewFromHex("512102cb560e47b1ae629416b4293256443cef4427cd5e5f233a8fd2a92f1912ece4a42103da972d5d07c3abd1c30938f11e9fa78b536ec833d5ba5fc0f2146aaa7d48e48352ae")

	p2pkhUtxo := &transaction.UTXO{
		TxID:          prevTx,
		Vout:          1,
		LockingScript: s,
		Satoshis:      111,
		Template:      tmpl,
	}
	tx.FromUTXOs(p2pkhUtxo)

	tx.Sign()
}
