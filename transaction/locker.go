package scripttemplate

import script "github.com/bitcoin-sv/go-sdk/script"

type Locker interface {
	LockingScript() *script.Script
	IsLockingScript(script *script.Script) bool
	Parse(script *script.Script) error
}
