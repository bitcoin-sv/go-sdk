package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/util"
)

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

// AddP2PKHInputsFromTx will add all Outputs of given previous transaction
// that match a specific public key to your transaction.
func (tx *Transaction) AddP2PKHInputsFromTx(pvsTx *Transaction, matchPK []byte) error {
	// Given that the prevTxID never changes, calculate it once up front.
	prevTxIDBytes := pvsTx.TxIDBytes()
	for i, utxo := range pvsTx.Outputs {
		utxoPkHASH160, err := utxo.LockingScript.PublicKeyHash()
		if err != nil {
			return err
		}

		if bytes.Equal(utxoPkHASH160, crypto.Hash160(matchPK)) {
			if err := tx.FromUTXOs(&UTXO{
				TxID:          prevTxIDBytes,
				Vout:          uint32(i),
				Satoshis:      utxo.Satoshis,
				LockingScript: utxo.LockingScript,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// From adds a new input to the transaction from the specified UTXO fields, using the default
// finalised sequence number (0xFFFFFFFF). If you want a different nSeq, change it manually
// afterwards.
func (tx *Transaction) From(prevTxID string, vout uint32, prevTxLockingScript string, satoshis uint64) error {
	pts, err := script.NewFromHex(prevTxLockingScript)
	if err != nil {
		return err
	}
	pti, err := hex.DecodeString(prevTxID)
	if err != nil {
		return err
	}

	return tx.FromUTXOs(&UTXO{
		TxID:          pti,
		Vout:          vout,
		LockingScript: pts,
		Satoshis:      satoshis,
	})
}

// FromUTXOs adds a new input to the transaction from the specified *bt.UTXO fields, using the default
// finalised sequence number (0xFFFFFFFF). If you want a different nSeq, change it manually
// afterwards.
func (tx *Transaction) FromUTXOs(utxos ...*UTXO) error {
	for _, utxo := range utxos {
		i := &TransactionInput{
			PreviousTxID:       utxo.TxID,
			PreviousTxOutIndex: utxo.Vout,
			PreviousTxSatoshis: utxo.Satoshis,
			PreviousTxScript:   utxo.LockingScript,
			SequenceNumber:     DefaultSequenceNumber, // use default finalised sequence number
		}

		tx.AddInput(i)
	}

	return nil
}
