package transaction

import (
	"context"
	"encoding/binary"

	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/util"
)

// UTXOGetterFunc is used for tx.Fund(...). It provides the amount of satoshis required
// for funding as `deficit`, and expects []*bt.UTXO to be returned containing
// utxos of which *bt.Input's can be built.
// If the returned []*bt.UTXO does not cover the deficit after fee recalculation, then
// this UTXOGetterFunc is called again, with the newly calculated deficit passed in.
//
// It is expected that bt.ErrNoUTXO will be returned once the utxo source is depleted.
type UTXOGetterFunc func(ctx context.Context, deficit uint64) ([]*UTXO, error)

// TotalInputSatoshis returns the total Satoshis inputted to the transaction.
func (tx *Transaction) TotalInputSatoshis() (total uint64) {
	for _, in := range tx.Inputs {
		total += in.PreviousTxSatoshis
	}
	return
}

func (tx *Transaction) AddInput(input *TransactionInput) {
	tx.Inputs = append(tx.Inputs, input)
}

func (tx *Transaction) AddInputFromTx(prevTx *Transaction, vout uint32) {
	i := &TransactionInput{
		PreviousTxOutIndex: vout,
		SequenceNumber:     DefaultSequenceNumber, // use default finalised sequence number
	}
	i.SetPreviousTx(prevTx)

	tx.Inputs = append(tx.Inputs, i)
}

// InputCount returns the number of transaction Inputs.
func (tx *Transaction) InputCount() int {
	return len(tx.Inputs)
}

// PreviousOutHash returns a byte slice of inputs outpoints, for creating a signature hash
func (tx *Transaction) PreviousOutHash() []byte {
	buf := make([]byte, 0)

	for _, in := range tx.Inputs {
		buf = append(buf, util.ReverseBytes(in.PreviousTxID)...)
		oi := make([]byte, 4)
		binary.LittleEndian.PutUint32(oi, in.PreviousTxOutIndex)
		buf = append(buf, oi...)
	}

	return crypto.Sha256d(buf)
}

// SequenceHash returns a byte slice of inputs SequenceNumber, for creating a signature hash
func (tx *Transaction) SequenceHash() []byte {
	buf := make([]byte, 0)

	for _, in := range tx.Inputs {
		oi := make([]byte, 4)
		binary.LittleEndian.PutUint32(oi, in.SequenceNumber)
		buf = append(buf, oi...)
	}

	return crypto.Sha256d(buf)
}
