package transaction

import (
	"encoding/binary"
	"encoding/hex"

	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/util"
)

// TotalInputSatoshis returns the total Satoshis inputted to the transaction.
func (tx *Transaction) TotalInputSatoshis() (total uint64) {
	for _, in := range tx.Inputs {
		prevSats := uint64(0)
		if in.SourceTxSatoshis() != nil {
			prevSats = *in.SourceTxSatoshis()
		}
		total += prevSats
	}
	return
}

func (tx *Transaction) AddInput(input *TransactionInput) {
	tx.Inputs = append(tx.Inputs, input)
}

func (tx *Transaction) AddInputWithOutput(input *TransactionInput, output *TransactionOutput) {
	input.SetPrevTxFromOutput(output)
	tx.Inputs = append(tx.Inputs, input)
}

func (tx *Transaction) AddInputFromTx(prevTx *Transaction, vout uint32,
	unlockingScriptTemplate UnlockingScriptTemplate) {
	i := &TransactionInput{
		SourceTXID:              prevTx.TxIDBytes(),
		SourceTxOutIndex:        vout,
		SourceTransaction:       prevTx,
		SequenceNumber:          DefaultSequenceNumber, // use default finalized sequence number
		UnlockingScriptTemplate: unlockingScriptTemplate,
	}

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
		buf = append(buf, util.ReverseBytes(in.SourceTXID)...)
		oi := make([]byte, 4)
		binary.LittleEndian.PutUint32(oi, in.SourceTxOutIndex)
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

// AddInputFrom adds a new input to the transaction from the specified UTXO fields, using the default
// finalized sequence number (0xFFFFFFFF). If you want a different nSeq, change it manually
// afterwards.
func (tx *Transaction) AddInputFrom(prevTxID string, vout uint32, prevTxLockingScript string,
	satoshis uint64, unlockingScriptTemplate UnlockingScriptTemplate) error {
	pts, err := script.NewFromHex(prevTxLockingScript)
	if err != nil {
		return err
	}
	pti, err := hex.DecodeString(prevTxID)
	if err != nil {
		return err
	}

	return tx.AddInputsFromUTXOs(&UTXO{
		TxID:                    pti,
		Vout:                    vout,
		LockingScript:           pts,
		Satoshis:                satoshis,
		UnlockingScriptTemplate: unlockingScriptTemplate,
	})
}

// AddInputsFromUTXOs adds a new input to the transaction from the specified *bt.UTXO fields, using the default
// finalized sequence number (0xFFFFFFFF). If you want a different nSeq, change it manually
// afterwards.
func (tx *Transaction) AddInputsFromUTXOs(utxos ...*UTXO) error {
	for _, utxo := range utxos {
		i := &TransactionInput{
			SourceTXID:              utxo.TxID,
			SourceTxOutIndex:        utxo.Vout,
			SequenceNumber:          DefaultSequenceNumber, // use default finalized sequence number
			UnlockingScriptTemplate: utxo.UnlockingScriptTemplate,
		}
		i.SetPrevTxFromOutput(&TransactionOutput{
			Satoshis:      utxo.Satoshis,
			LockingScript: utxo.LockingScript,
		})

		tx.AddInput(i)
	}

	return nil
}
