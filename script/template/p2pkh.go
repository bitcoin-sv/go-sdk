package template

import (
	"github.com/bitcoin-sv/go-sdk/script"
)

type P2PKH struct {
	PubKeyHash []byte
}

func (p *P2PKH) lock() (s *script.Script, err error) {
	return script.NewScriptFromScriptOps([]*script.ScriptOp{
		&script.ScriptOp{OpCode: script.OpDUP},
		&script.ScriptOp{OpCode: script.OpHASH160},
		&script.ScriptOp{OpCode: script.OpDATA20, Data: p.PubKeyHash},
		&script.ScriptOp{OpCode: script.OpEQUALVERIFY},
		&script.ScriptOp{OpCode: script.OpCHECKSIG},
	})
}
