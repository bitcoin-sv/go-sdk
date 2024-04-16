package locker

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec"
)

type P2PKH struct {
	PubKey *ec.PublicKey
}

func (p P2PKH) LockingScript() (*bscript.Script, error) {
	return bscript.NewP2PKHFromPubKeyBytes(p.PubKey.SerialiseCompressed())
}

func (p P2PKH) IsLockingScript(script *bscript.Script) bool {
	return script.IsP2PKH()
}
