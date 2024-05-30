package template

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type MultisigTemplate struct {
	PubKeys      []*ec.PublicKey
	RequiredSigs int
	privKeys     []*ec.PrivateKey
}

func NewMultisigTemplateFromPrivKeys(privKeys []*ec.PrivateKey, requiredSigs int) *MultisigTemplate {
	pubKeys := make([]*ec.PublicKey, len(privKeys))
	for i, privKey := range privKeys {
		pubKeys[i] = privKey.PubKey()
	}
	return &MultisigTemplate{
		PubKeys:      pubKeys,
		RequiredSigs: requiredSigs,
		privKeys:     privKeys,
	}
}

func (m *MultisigTemplate) IsLockingScript(script *script.Script) bool {
	return script.IsMultiSigOut()
}

func (m *MultisigTemplate) IsUnlockingScript(s *script.Script) bool {
	pos := 0
	if op, err := s.ReadOp(&pos); err != nil {
		return false
	} else if op.OpCode != script.Op0 {
		return false
	}
	for {
		if op, err := s.ReadOp(&pos); err == script.ErrScriptIndexOutOfRange {
			return true
		} else if err != nil {
			return false
		} else if _, err := ec.FromDER(op.Data); err != nil {
			return false
		}
	}
}

func (m *MultisigTemplate) Lock() (*script.Script, error) {
	if m.RequiredSigs > 16 || len(m.PubKeys) > 16 {
		return nil, ErrTooManySignatures
	}
	s := &script.Script{}
	s.AppendOpcodes(80 + uint8(m.RequiredSigs))
	for _, pubKey := range m.PubKeys {
		s.AppendPushData(pubKey.SerialiseCompressed())
	}
	s.AppendOpcodes(80 + uint8(len(m.PubKeys)))
	s.AppendOpcodes(script.OpCHECKMULTISIG)

	return s, nil
}

func (m *MultisigTemplate) Sign(tx *transaction.Transaction, params transaction.UnlockParams) (*script.Script, error) {
	if len(m.privKeys) == 0 {
		return nil, ErrNoPrivateKey
	}
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}

	if tx.Inputs[params.InputIdx].PreviousTxScript == nil {
		return nil, transaction.ErrEmptyPreviousTxScript
	}
	s := tx.Inputs[params.InputIdx].PreviousTxScript

	if !m.IsLockingScript(s) {
		return nil, ErrBadScript
	}

	uscript := &script.Script{}
	pos := 1
	for i := 0; i < 16; i++ {
		op, err := s.ReadOp(&pos)
		if err != nil {
			return nil, err
		}

		if len(op.Data) == 0 && op.OpCode != script.OpCHECKMULTISIG {
			return nil, ErrBadScript
		}

		if op.OpCode == script.OpCHECKMULTISIG {
			return uscript, nil
		}

		pubKey, err := ec.ParsePubKey(op.Data)
		if err != nil {
			return nil, err
		}
		for j, p := range m.PubKeys {
			if pubKey.IsEqual(p) {
				sh, err := tx.CalcInputSignatureHash(params.InputIdx, params.SigHashFlags)
				if err != nil {
					return nil, err
				}
				sig, err := m.privKeys[j].Sign(sh)
				if err != nil {
					return nil, err
				}

				uscript.AppendPushData(sig.Serialise())
				break
			}
		}
	}

	return uscript, nil
}

func (p *MultisigTemplate) EstimateSize(_ *transaction.Transaction, inputIndex uint32) int {
	return 34 * p.RequiredSigs
}
