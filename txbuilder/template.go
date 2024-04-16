package txbuilder

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

type Template interface {
	LockingScript() *bscript.Script
	IsLockingScript() bool
	IsUnlockingScript() bool
	UnlockingScript(tx *transaction.Tx, up transaction.UnlockerParams) (uscript *bscript.Script, err error)
	EstimateSize() uint32
}
