package unlocker

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec"
	"github.com/bitcoin-sv/go-sdk/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/locker"
)

type Multisig struct {
	PrivateKey *ec.PrivateKey
}

func (l *Multisig) UnlockingScript(ctx context.Context, tx *transaction.Transaction, params transaction.UnlockerParams) (*bscript.Script, error) {
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

	var uscript *bscript.Script
	if len(*input.UnlockingScript) > 0 {
		uscript = input.UnlockingScript
		pos := 0
		ops := []*bscript.ScriptOp{}
		last0 := 0
		for i := range lock.RequiredSigs {
			op, err := uscript.ReadOp(&pos)
			if err != nil {
				if err == bscript.ErrScriptIndexOutOfRange {
					break
				}
				return nil, err
			}
			if op.OpCode == bscript.Op0 {
				last0 = i
			}
			ops = append(ops, op)
		}
		ops[last0] = &bscript.ScriptOp{OpCode: byte(len(signature)), Data: signature}
		return bscript.NewScriptFromScriptOps(ops)

	} else {
		uscript = &bscript.Script{}
		// 1 OP_0 for each required signature.
		for range lock.RequiredSigs {
			uscript.AppendOpcodes(bscript.Op0)
		}
		// Apply signature.
		// This leaves an extra push data on the stack to resolve OP_CHECKMULTISIG bug
		uscript.AppendPushData(signature)
	}
	return uscript, nil
}
