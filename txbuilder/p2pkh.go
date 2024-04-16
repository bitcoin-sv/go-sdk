package txbuilder

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

type P2PKH struct {
	bscript.Script
	PubKey *ec.PublicKey
}

func (p *P2PKH) LockingScript() (*bscript.Script, error) {
	return bscript.NewP2PKHFromPubKeyBytes(p.PubKey.SerialiseCompressed())
}

func (p *P2PKH) IsLockingScript() bool {
	return p.IsP2PKH()
}

func (p *P2PKH) IsUnlockingScript() bool {
	return true
}

func (p *P2PKH) UnlockingScript(tx transaction.Tx, params transaction.UnlockerParams) (*bscript.Script, error) {
	return &bscript.Script{}, nil
}

func (p *P2PKH) EstimateSize() uint32 {
	return 25
}
