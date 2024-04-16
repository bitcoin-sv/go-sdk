package transaction

import "github.com/bitcoin-sv/go-sdk/bscript"

type Locker interface {
	LockingScript() *bscript.Script
	IsLockingScript(script *bscript.Script) bool
	Parse(script *bscript.Script) error
}
