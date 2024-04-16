package txbuilder

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

type UnlockerParams struct {
	// InputIdx the input to be unlocked. [DEFAULT 0]
	InputIdx uint32
	// SigHashFlags the be applied [DEFAULT ALL|FORKID]
	SigHashFlags sighash.Flag
}

type Template interface {
	LockingScript() *bscript.Script
	IsLockingScript() bool
	IsUnlockingScript() bool
	UnlockingScript(tx *transaction.Tx, up UnlockerParams) (uscript *bscript.Script, err error)
	EstimateSize() uint32
}
