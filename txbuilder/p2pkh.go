package txbuilder

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

type P2PKH struct {
	bscript.Script
	PubKey []byte
}

func (p *P2PKH) Lock() (*bscript.Script, error) {
	return bscript.NewP2PKHFromPubKeyBytes(p.PubKey)
}

func (p *P2PKH) IsLockingScript() bool {
	return p.IsP2PKH()
}

func (p *P2PKH) IsUnlockingScript() bool {
	return true
}

func (p *P2PKH) Sign(tx transaction.Tx, params transaction.UnlockerParams) (*bscript.Script, error) {
	return &bscript.Script{}, nil
}

func (p *P2PKH) EstimateSize() uint32 {
	return 25
}
