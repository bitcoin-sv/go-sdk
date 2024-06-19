package transaction

import (
	"github.com/bitcoin-sv/go-sdk/script"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type UnlockParams struct {
	// InputIdx the input to be unlocked. [DEFAULT 0]
	InputIdx uint32
	// SigHashFlags the be applied [DEFAULT ALL|FORKID]
	SigHashFlags     sighash.Flag
	CodeSeparatorIdx uint32
}

type Unlocker interface {
	Unlock(tx *Transaction, unlockParams *UnlockParams) (*script.Script, error)
	EstimateLength(tx *Transaction, inputIndex uint32) uint32
}

type ScriptTemplate interface {
	Lock(params ...any) (*script.Script, error)
	Unlocker(params ...any) (Unlocker, error)
}
