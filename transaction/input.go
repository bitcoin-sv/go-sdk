package transaction

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/util"
	"github.com/pkg/errors"
)

/*
Field	                     Description                                                   Size
--------------------------------------------------------------------------------------------------------
Previous Transaction hash  doubled SHA256-hashed of a (previous) to-be-used transaction	 32 bytes
Previous Txout-index       non-negative integer indexing an output of the to-be-used      4 bytes
                           transaction
Txin-script length         non-negative integer VI = VarInt                               1-9 bytes
Txin-script / scriptSig	   Script	                                                        <in-script length>-many bytes
sequence_no	               normally 0xFFFFFFFF; irrelevant unless transaction's           4 bytes
                           lock_time is > 0
*/

// DefaultSequenceNumber is the default starting sequence number
const DefaultSequenceNumber uint32 = 0xFFFFFFFF

// TransactionInput is a representation of a transaction input
//
// DO NOT CHANGE ORDER - Optimised for memory via maligned
type TransactionInput struct {
	previousTx         *Transaction
	PreviousTxID       []byte
	PreviousTxSatoshis uint64
	PreviousTxScript   *script.Script
	UnlockingScript    *script.Script
	PreviousTxOutIndex uint32
	SequenceNumber     uint32
	Unlocker           *Unlocker
}

func (i *TransactionInput) PreviousTx() *Transaction {
	return i.previousTx
}

func (i *TransactionInput) SetPreviousTx(tx *Transaction) {
	i.PreviousTxID = tx.TxIDBytes()
	i.PreviousTxScript = tx.Outputs[i.PreviousTxOutIndex].LockingScript
	i.PreviousTxSatoshis = tx.Outputs[i.PreviousTxOutIndex].Satoshis
	i.previousTx = tx
}

// ReadFrom reads from the `io.Reader` into the `bt.Input`.
func (i *TransactionInput) ReadFrom(r io.Reader) (int64, error) {
	return i.readFrom(r, false)
}

// ReadFromExtended reads the `io.Reader` into the `bt.Input` when the reader is
// consuming an extended format transaction.
func (i *TransactionInput) ReadFromExtended(r io.Reader) (int64, error) {
	return i.readFrom(r, true)
}

func (i *TransactionInput) readFrom(r io.Reader, extended bool) (int64, error) {
	*i = TransactionInput{}
	var bytesRead int64

	previousTxID := make([]byte, 32)
	n, err := io.ReadFull(r, previousTxID)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, errors.Wrapf(err, "previousTxID(32): got %d bytes", n)
	}

	prevIndex := make([]byte, 4)
	n, err = io.ReadFull(r, prevIndex)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, errors.Wrapf(err, "previousTxID(4): got %d bytes", n)
	}

	var l VarInt
	n64, err := l.ReadFrom(r)
	bytesRead += n64
	if err != nil {
		return bytesRead, err
	}

	scriptBytes := make([]byte, l)
	n, err = io.ReadFull(r, scriptBytes)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, errors.Wrapf(err, "script(%d): got %d bytes", l, n)
	}

	sequence := make([]byte, 4)
	n, err = io.ReadFull(r, sequence)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, errors.Wrapf(err, "sequence(4): got %d bytes", n)
	}

	i.PreviousTxID = util.ReverseBytes(previousTxID)
	i.PreviousTxOutIndex = binary.LittleEndian.Uint32(prevIndex)
	i.UnlockingScript = script.NewFromBytes(scriptBytes)
	i.SequenceNumber = binary.LittleEndian.Uint32(sequence)

	if extended {
		prevSatoshis := make([]byte, 8)
		var prevTxLockingScript script.Script

		n, err = io.ReadFull(r, prevSatoshis)
		bytesRead += int64(n)
		if err != nil {
			return bytesRead, errors.Wrapf(err, "prevSatoshis(8): got %d bytes", n)
		}

		// Read in the prevTxLockingScript
		var scriptLen VarInt
		n64, err := scriptLen.ReadFrom(r)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}

		scriptBytes := make([]byte, scriptLen)
		n, err := io.ReadFull(r, scriptBytes)
		bytesRead += int64(n)
		if err != nil {
			return bytesRead, errors.Wrapf(err, "script(%d): got %d bytes", scriptLen.Length(), n)
		}

		prevTxLockingScript = *script.NewFromBytes(scriptBytes)

		i.PreviousTxSatoshis = binary.LittleEndian.Uint64(prevSatoshis)
		i.PreviousTxScript = script.NewFromBytes(prevTxLockingScript)
	}

	return bytesRead, nil
}

// PreviousTxIDStr returns the Previous TxID as a hex string.
func (i *TransactionInput) PreviousTxIDStr() string {
	return hex.EncodeToString(i.PreviousTxID)
}

// String implements the Stringer interface and returns a string
// representation of a transaction input.
func (i *TransactionInput) String() string {
	return fmt.Sprintf(
		`prevTxHash:   %s
prevOutIndex: %d
scriptLen:    %d
script:       %s
sequence:     %x
`,
		hex.EncodeToString(i.PreviousTxID),
		i.PreviousTxOutIndex,
		len(*i.UnlockingScript),
		i.UnlockingScript,
		i.SequenceNumber,
	)
}

// Bytes encodes the Input into a hex byte array.
func (i *TransactionInput) Bytes(clear bool) []byte {
	h := make([]byte, 0)

	h = append(h, util.ReverseBytes(i.PreviousTxID)...)
	h = append(h, util.LittleEndianBytes(i.PreviousTxOutIndex, 4)...)
	if clear {
		h = append(h, 0x00)
	} else {
		if i.UnlockingScript == nil {
			h = append(h, VarInt(0).Bytes()...)
		} else {
			h = append(h, VarInt(uint64(len(*i.UnlockingScript))).Bytes()...)
			h = append(h, *i.UnlockingScript...)
		}
	}

	return append(h, util.LittleEndianBytes(i.SequenceNumber, 4)...)
}
