package transaction

import (
	"github.com/bitcoin-sv/go-sdk/script"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type UnlockParams struct {
	// InputIdx the input to be unlocked. [DEFAULT 0]
	InputIdx uint32
	// SigHashFlags the be applied [DEFAULT ALL|FORKID]
	SigHashFlags   sighash.Flag
	LockingScript  *script.Script
	SourceSatoshis uint64
}

type ScriptTemplate interface {
	IsLockingScript(s *script.Script) bool
	IsUnlockingScript(s *script.Script) bool
	Lock() (*script.Script, error)
	Sign(tx *Transaction, up UnlockParams) (uscript *script.Script, err error)
	EstimateSize(tx *Transaction, inputIndex uint32) int
}
