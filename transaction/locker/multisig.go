package locker

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/bscript"
	"github.com/bitcoin-sv/go-sdk/ec"
)

type Multisig struct {
	PubKeys      []*ec.PublicKey
	RequiredSigs int
}

func (m Multisig) LockingScript() (*bscript.Script, error) {
	if m.RequiredSigs > 16 || len(m.PubKeys) > 16 {
		return nil, fmt.Errorf("Multisig: Sigs must be less than 16")
	}
	script := &bscript.Script{}
	script.AppendOpcodes(80 + uint8(m.RequiredSigs))
	for _, pubKey := range m.PubKeys {
		script.AppendPushData(pubKey.SerialiseCompressed())
	}
	script.AppendOpcodes(80 + uint8(len(m.PubKeys)))
	script.AppendOpcodes(bscript.OpCHECKMULTISIG)

	return script, nil
}

func (m Multisig) IsLockingScript(script *bscript.Script) bool {
	return script.IsMultiSigOut()
}

func (m *Multisig) Parse(script *bscript.Script) error {
	if !m.IsLockingScript(script) {
		return fmt.Errorf("Multisig: script is not multisig")
	}

	pos := 0
	op, err := script.ReadOp(&pos)
	if err != nil {
		return err
	}

	if op.OpCode < 0x50 || op.OpCode > 0x5f {
		return fmt.Errorf("Multisig: script is not multisig")
	}
	m.RequiredSigs = int(op.OpCode - 0x50)

	for i := 0; i < 16; i++ {
		op, err = script.ReadOp(&pos)
		if err != nil {
			return err
		}

		if len(op.Data) == 0 && op.OpCode != bscript.OpCHECKMULTISIG {
			return fmt.Errorf("Multisig: script is not multisig")
		}

		if op.OpCode == bscript.OpCHECKMULTISIG {
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
