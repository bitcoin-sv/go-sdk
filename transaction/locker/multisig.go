package locker

import (
	"fmt"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	script "github.com/bitcoin-sv/go-sdk/script"
)

type Multisig struct {
	PubKeys      []*ec.PublicKey
	RequiredSigs int
}

func (m Multisig) LockingScript() (*script.Script, error) {
	if m.RequiredSigs > 16 || len(m.PubKeys) > 16 {
		return nil, fmt.Errorf("Multisig: Sigs must be less than 16")
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

func (m Multisig) IsLockingScript(script *script.Script) bool {
	return script.IsMultiSigOut()
}

func (m *Multisig) Parse(s *script.Script) error {
	if !m.IsLockingScript(s) {
		return fmt.Errorf("Multisig: script is not multisig")
	}

	pos := 0
	op, err := s.ReadOp(&pos)
	if err != nil {
		return err
	}

	if op.OpCode < 0x50 || op.OpCode > 0x5f {
		return fmt.Errorf("Multisig: script is not multisig")
	}
	m.RequiredSigs = int(op.OpCode - 0x50)

	for i := 0; i < 16; i++ {
		op, err = s.ReadOp(&pos)
		if err != nil {
			return err
		}

		if len(op.Data) == 0 && op.OpCode != script.OpCHECKMULTISIG {
			return fmt.Errorf("Multisig: script is not multisig")
		}

		if op.OpCode == script.OpCHECKMULTISIG {
			return nil
		}

		pubKey, err := ec.ParsePubKey(op.Data)
		if err != nil {
			return err
		}
		m.PubKeys = append(m.PubKeys, pubKey)
	}

	return nil
}
