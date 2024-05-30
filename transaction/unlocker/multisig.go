package unlocker

import (
	"context"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/locker"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type Multisig struct {
	PrivateKey *ec.PrivateKey
}

func (l *Multisig) UnlockingScript(ctx context.Context, tx *transaction.Transaction, params transaction.UnlockerParams) (*script.Script, error) {
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}
	input := tx.Inputs[params.InputIdx]

	if input.PreviousTxScript == nil {
		return nil, transaction.ErrEmptyPreviousTxScript
	}

	lock := &locker.Multisig{}
	err := lock.Parse(input.PreviousTxScript)
	if err != nil {
		return nil, err
	}

	sh, err := tx.CalcInputSignatureHash(params.InputIdx, params.SigHashFlags)
	if err != nil {
		return nil, err
	}

	sig, err := l.PrivateKey.Sign(sh)
	if err != nil {
		return nil, err
	}

	signature := sig.Serialise()

	var uscript *script.Script
	if len(*input.UnlockingScript) > 0 {
		uscript = input.UnlockingScript
		pos := 0
		ops := []*script.ScriptOp{}
		last0 := 0
		for i := range lock.RequiredSigs {
			op, err := uscript.ReadOp(&pos)
			if err != nil {
				if err == script.ErrScriptIndexOutOfRange {
					break
				}
				return nil, err
			}
			if op.OpCode == script.Op0 {
				last0 = i
			}
			ops = append(ops, op)
		}
		ops[last0] = &script.ScriptOp{OpCode: byte(len(signature)), Data: signature}
		return script.NewScriptFromScriptOps(ops)

	} else {
		uscript = &script.Script{}
		// 1 OP_0 for each required signature.
		for range lock.RequiredSigs {
			uscript.AppendOpcodes(script.Op0)
		}
		// Apply signature.
		// This leaves an extra push data on the stack to resolve OP_CHECKMULTISIG bug
		uscript.AppendPushData(signature)
	}
	return uscript, nil
}
