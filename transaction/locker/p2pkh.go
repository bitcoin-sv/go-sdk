package locker

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
)

type P2PKH struct {
	PubKey *ec.PublicKey
}

func (p P2PKH) LockingScript() (*script.Script, error) {
	return script.NewP2PKHFromPubKeyBytes(p.PubKey.SerialiseCompressed())
}

func (p P2PKH) IsLockingScript(script *script.Script) bool {
	return script.IsP2PKH()
}
