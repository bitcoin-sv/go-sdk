package transaction

import (
	"encoding/hex"

	bscript "github.com/bitcoin-sv/go-sdk/script"
)

// UTXO an unspent transaction output, used for creating inputs
type UTXO struct {
	TxID           []byte          `json:"txid"`
	Vout           uint32          `json:"vout"`
	LockingScript  *bscript.Script `json:"locking_script"`
	Satoshis       uint64          `json:"satoshis"`
	SequenceNumber uint32          `json:"sequence_number"`
	Unlocker       *Unlocker       `json:"-"`
}

// UTXOs a collection of *bt.UTXO.
type UTXOs []*UTXO

// NodeJSON returns a wrapped *bt.UTXO for marshalling/unmarshalling into a node utxo format.
//
// Marshalling usage example:
//
//	bb, err := json.Marshal(utxo.NodeJSON())
//
// Unmarshalling usage example:
//
//	utxo := &bt.UTXO{}
//	if err := json.Unmarshal(bb, utxo.NodeJSON()); err != nil {}
func (u *UTXO) NodeJSON() interface{} {
	return &nodeUTXOWrapper{UTXO: u}
}

// NodeJSON returns a wrapped bt.UTXOs for marshalling/unmarshalling into a node utxo format.
//
// Marshalling usage example:
//
//	bb, err := json.Marshal(utxos.NodeJSON())
//
// Unmarshalling usage example:
//
//	var txs bt.UTXOs
//	if err := json.Unmarshal(bb, utxos.NodeJSON()); err != nil {}
func (u *UTXOs) NodeJSON() interface{} {
	return (*nodeUTXOsWrapper)(u)
}

// TxIDStr return the tx id as a string.
func (u *UTXO) TxIDStr() string {
	return hex.EncodeToString(u.TxID)
}

// LockingScriptHex retur nthe locking script in hex format.
func (u *UTXO) LockingScriptHex() string {
	return u.LockingScript.String()
}

// NewUTXO creates a new UTXO.
func NewUTXO(prevTxID string, vout uint32, prevTxLockingScript string, satoshis uint64) (*UTXO, error) {
	pts, err := bscript.NewFromHex(prevTxLockingScript)
	if err != nil {
		return nil, err
	}
	pti, err := hex.DecodeString(prevTxID)
	if err != nil {
		return nil, err
	}

	return &UTXO{
		TxID:          pti,
		Vout:          vout,
		LockingScript: pts,
		Satoshis:      satoshis,
	}, nil
}
