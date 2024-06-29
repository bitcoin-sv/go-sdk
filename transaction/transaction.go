package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"slices"

	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/util"
)

type Transaction struct {
	Version    uint32               `json:"version"`
	Inputs     []*TransactionInput  `json:"inputs"`
	Outputs    []*TransactionOutput `json:"outputs"`
	LockTime   uint32               `json:"locktime"`
	MerklePath *MerklePath          `json:"merklePath"`
}

// Transactions a collection of *bt.Tx.
type Transactions []*Transaction

// NewTransaction creates a new transaction object with default values.
func NewTransaction() *Transaction {
	return &Transaction{Version: 1, LockTime: 0, Inputs: make([]*TransactionInput, 0)}
}

// NewTransactionFromHex takes a toBytesHelper string representation of a bitcoin transaction
// and returns a Tx object.
func NewTransactionFromHex(str string) (*Transaction, error) {
	bb, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return NewTransactionFromBytes(bb)
}

// NewTransactionFromBytes takes an array of bytes, constructs a Tx and returns it.
// This function assumes that the byte slice contains exactly 1 transaction.
func NewTransactionFromBytes(b []byte) (*Transaction, error) {
	tx, used, err := NewTransactionFromStream(b)
	if err != nil {
		return nil, err
	}

	if used != len(b) {
		return nil, ErrNLockTimeLength
	}

	return tx, nil
}

// NewTransactionFromStream takes an array of bytes and constructs a Tx from it, returning the Tx and the bytes used.
// Despite the name, this is not actually reading a stream in the true sense: it is a byte slice that contains
// many transactions one after another.
func NewTransactionFromStream(b []byte) (*Transaction, int, error) {
	tx := Transaction{}

	bytesRead, err := tx.ReadFrom(bytes.NewReader(b))

	return &tx, int(bytesRead), err
}

// ReadFrom reads from the `io.Reader` into the `bt.Tx`.
func (tx *Transaction) ReadFrom(r io.Reader) (int64, error) {
	*tx = Transaction{}
	var bytesRead int64

	// Define n64 and err here to avoid linter complaining about shadowing variables.
	var n64 int64
	var err error

	version := make([]byte, 4)
	n, err := io.ReadFull(r, version)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, err
	}

	tx.Version = binary.LittleEndian.Uint32(version)

	extended := false

	var inputCount VarInt

	n64, err = inputCount.ReadFrom(r)
	bytesRead += n64
	if err != nil {
		return bytesRead, err
	}

	var outputCount VarInt
	locktime := make([]byte, 4)

	// ----------------------------------------------------------------------------------
	// If the inputCount is 0, we may be parsing an incomplete transaction, or we may be
	// both of these cases without needing to rewind (peek) the incoming stream of bytes.
	// ----------------------------------------------------------------------------------
	if inputCount == 0 {
		n64, err = outputCount.ReadFrom(r)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}

		if outputCount == 0 {
			// Read in lock time
			n, err = io.ReadFull(r, locktime)
			bytesRead += int64(n)
			if err != nil {
				return bytesRead, err
			}

			if binary.BigEndian.Uint32(locktime) != 0xEF {
				tx.LockTime = binary.LittleEndian.Uint32(locktime)
				return bytesRead, nil
			}

			extended = true

			n64, err = inputCount.ReadFrom(r)
			bytesRead += n64
			if err != nil {
				return bytesRead, err
			}
		}
	}
	// ----------------------------------------------------------------------------------
	// If we have not returned from the previous block, we will have detected a sane
	// transaction and we will know if it is extended format or not.
	// We can now proceed with reading the rest of the transaction.
	// ----------------------------------------------------------------------------------

	// create Inputs
	for i := uint64(0); i < uint64(inputCount); i++ {
		input := &TransactionInput{}
		n64, err = input.readFrom(r, extended)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}
		tx.Inputs = append(tx.Inputs, input)
	}

	if inputCount > 0 || extended {
		// Re-read the actual output count...
		n64, err = outputCount.ReadFrom(r)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}
	}

	for i := uint64(0); i < uint64(outputCount); i++ {
		output := new(TransactionOutput)
		n64, err = output.ReadFrom(r)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}

		tx.Outputs = append(tx.Outputs, output)
	}

	n, err = io.ReadFull(r, locktime)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, err
	}
	tx.LockTime = binary.LittleEndian.Uint32(locktime)

	return bytesRead, nil
}

// ReadFrom txs from a block in a `bt.Txs`. This assumes a preceding varint detailing
// the total number of txs that the reader will provide.
func (tt *Transactions) ReadFrom(r io.Reader) (int64, error) {
	var bytesRead int64

	var txCount VarInt
	n, err := txCount.ReadFrom(r)
	bytesRead += n
	if err != nil {
		return bytesRead, err
	}

	*tt = make([]*Transaction, txCount)

	for i := uint64(0); i < uint64(txCount); i++ {
		tx := new(Transaction)
		n, err := tx.ReadFrom(r)
		bytesRead += n
		if err != nil {
			return bytesRead, err
		}

		(*tt)[i] = tx
	}

	return bytesRead, nil
}

// HasDataOutputs returns true if the transaction has
// at least one data (OP_RETURN) output in it.
func (tx *Transaction) HasDataOutputs() bool {
	for _, out := range tx.Outputs {
		if out.LockingScript.IsData() {
			return true
		}
	}
	return false
}

// InputIdx will return the input at the specified index.
//
// This will consume an overflow error and simply return nil if the input
// isn't found at the index.
func (tx *Transaction) InputIdx(i int) *TransactionInput {
	if i > tx.InputCount()-1 {
		return nil
	}
	return tx.Inputs[i]
}

// OutputIdx will return the output at the specified index.
//
// This will consume an overflow error and simply return nil if the output
// isn't found at the index.
func (tx *Transaction) OutputIdx(i int) *TransactionOutput {
	if i > tx.OutputCount()-1 {
		return nil
	}
	return tx.Outputs[i]
}

// IsCoinbase determines if this transaction is a coinbase by
// checking if the tx input is a standard coinbase input.
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.Inputs) != 1 {
		return false
	}

	cbi := make([]byte, 32)

	if !bytes.Equal(tx.Inputs[0].SourceTXID, cbi) {
		return false
	}

	if tx.Inputs[0].SourceTxOutIndex == DefaultSequenceNumber || tx.Inputs[0].SequenceNumber == DefaultSequenceNumber {
		return true
	}

	return false
}

// TxIDBytes returns the transaction ID of the transaction as bytes
// (which is also the transaction hash).
func (tx *Transaction) TxIDBytes() []byte {
	return util.ReverseBytes(crypto.Sha256d(tx.Bytes()))
}

// TxID returns the transaction ID of the transaction
// (which is also the transaction hash).
func (tx *Transaction) TxID() string {
	return hex.EncodeToString(util.ReverseBytes(crypto.Sha256d(tx.Bytes())))
}

// String encodes the transaction into a hex string.
func (tx *Transaction) String() string {
	return hex.EncodeToString(tx.Bytes())
}

// IsValidTxID will check that the txid bytes are valid.
//
// A txid should be of 32 bytes length.
func IsValidTxID(txid []byte) bool {
	return len(txid) == 32
}

// Bytes encodes the transaction into a byte array.
// See https://chainquery.com/bitcoin-cli/decoderawtransaction
func (tx *Transaction) Bytes() []byte {
	return tx.toBytesHelper(0, nil, false)
}

// EF outputs the transaction into a byte array in extended format
// (with PreviousTxSatoshis and PreviousTXScript included)
func (tx *Transaction) EF() ([]byte, error) {
	for _, in := range tx.Inputs {
		if in.SourceTransaction == nil {
			return nil, ErrEmptyPreviousTx
		}
	}
	return tx.toBytesHelper(0, nil, true), nil
}

func (tx *Transaction) EFHex() (string, error) {
	ef, err := tx.EF()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ef), nil
}

// BytesWithClearedInputs encodes the transaction into a byte array but clears its Inputs first.
// This is used when signing transactions.
func (tx *Transaction) BytesWithClearedInputs(index int, lockingScript []byte) []byte {
	return tx.toBytesHelper(index, lockingScript, false)
}

// Clone returns a clone of the tx
func (tx *Transaction) Clone() *Transaction {
	// Ignore err as byte slice passed in is created from valid tx
	clone, err := NewTransactionFromBytes(tx.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	for i, input := range tx.Inputs {
		clone.Inputs[i].SourceTransaction = input.SourceTransaction
	}

	return clone
}

func (tx *Transaction) toBytesHelper(index int, lockingScript []byte, extended bool) []byte {
	h := make([]byte, 0)

	h = append(h, util.LittleEndianBytes(tx.Version, 4)...)

	if extended {
		h = append(h, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xEF}...)
	}

	h = append(h, VarInt(uint64(len(tx.Inputs))).Bytes()...)

	for i, in := range tx.Inputs {
		s := in.Bytes(lockingScript != nil)
		if i == index && lockingScript != nil {
			h = append(h, VarInt(uint64(len(lockingScript))).Bytes()...)
			h = append(h, lockingScript...)
		} else {
			h = append(h, s...)
		}

		if extended {
			b := make([]byte, 8)
			prevSats := uint64(0)
			if in.SourceTxSatoshis() != nil {
				prevSats = *in.SourceTxSatoshis()
			}
			binary.LittleEndian.PutUint64(b, prevSats)
			h = append(h, b...)

			prevScript := in.SourceTxScript()
			if prevScript != nil {
				l := uint64(len(*prevScript))
				h = append(h, VarInt(l).Bytes()...)
				h = append(h, *prevScript...)
			} else {
				h = append(h, 0x00) // The length of the script is zero
			}
		}
	}

	h = append(h, VarInt(uint64(len(tx.Outputs))).Bytes()...)
	for _, out := range tx.Outputs {
		h = append(h, out.Bytes()...)
	}

	lt := make([]byte, 4)
	binary.LittleEndian.PutUint32(lt, tx.LockTime)

	return append(h, lt...)
}

// Size will return the size of tx in bytes.
func (tx *Transaction) Size() int {
	return len(tx.Bytes())
}

func (tx *Transaction) AddMerkleProof(bump *MerklePath) error {
	if !slices.ContainsFunc(bump.Path[0], func(v *PathElement) bool {
		return bytes.Equal(v.Hash, tx.TxIDBytes())
	}) {
		return ErrBadMerkleProof
	}
	tx.MerklePath = bump
	return nil
}

func (tx *Transaction) Sign() error {
	for vin, i := range tx.Inputs {
		if i.UnlockingScriptTemplate != nil {
			unlock, err := i.UnlockingScriptTemplate.Sign(tx, uint32(vin))
			if err != nil {
				return err
			}
			i.UnlockingScript = unlock
		}
	}
	return nil
}
