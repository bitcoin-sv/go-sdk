package txbuilder

import (
	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

type Multisig struct {
	bscript.Script
	RequiredSigs int
	PubKeys      []*bscript.Script
}

func (m *Multisig) Lock() *bscript.Script {
	return &bscript.Script{}
}

func (m *Multisig) IsLockingScript() bool {
	return true
}

func (m *Multisig) IsUnlockingScript() bool {
	return false
}

func (m *Multisig) Sign(tx transaction.Tx, params transaction.UnlockerParams) (*bscript.Script, error) {
	return &bscript.Script{}, nil
}

func (m *Multisig) EstimateSize() uint32 {
	return 0
}
