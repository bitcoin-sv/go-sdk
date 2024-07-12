package transaction

import (
	"encoding/hex"

	script "github.com/bitcoin-sv/go-sdk/script"
)

// UTXO an unspent transaction output, used for creating inputs
type UTXO struct {
	TxID                    []byte                  `json:"txid"`
	Vout                    uint32                  `json:"vout"`
	LockingScript           *script.Script          `json:"locking_script"`
	Satoshis                uint64                  `json:"satoshis"`
	UnlockingScriptTemplate UnlockingScriptTemplate `json:"-"`
}

// UTXOs a collection of *bt.UTXO.
type UTXOs []*UTXO

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
	pts, err := script.NewFromHex(prevTxLockingScript)
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
