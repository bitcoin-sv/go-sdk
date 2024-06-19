package transaction

import (
	"github.com/bitcoin-sv/go-sdk/script"
)

type Unlocker interface {
	Unlock(tx *Transaction, inputIndex uint32) (*script.Script, error)
	EstimateLength(tx *Transaction, inputIndex uint32) uint32
}

// type ScriptTemplate interface {
// 	Lock(lockParams ...interface{}) (*script.Script, error)
// 	Unlocker(unlockParams ...interface{}) (Unlocker, error)
// }
