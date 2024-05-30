package template

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	hash "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/script"
)

type P2PKHTemplate struct {
	PKHash     []byte
	PrivateKey *ec.PrivateKey
}

func NewP2PKHTemplateFromAddress(address script.Address) (*P2PKHTemplate, error) {
	return &P2PKHTemplate{
		PKHash: address.PublicKeyHash,
	}, nil
}

func NewP2PKHTemplateFromPubKey(pubKey []byte) (*P2PKHTemplate, error) {
	return &P2PKHTemplate{
		PKHash: hash.Hash160(pubKey),
	}, nil
}

func NewP2PKHTemplateFromPubKeyEC(pubKey *ec.PublicKey) (*P2PKHTemplate, error) {
	return &P2PKHTemplate{
		PKHash: hash.Hash160(pubKey.SerialiseCompressed()),
	}, nil
}

func (p P2PKHTemplate) IsLockingScript(s *script.Script) bool {
	b := []byte(*s)
	return len(b) == 25 &&
		b[0] == script.OpDUP &&
		b[1] == script.OpHASH160 &&
		b[2] == script.OpDATA20 &&
		b[23] == script.OpEQUALVERIFY &&
		b[24] == script.OpCHECKSIG
}

func (p P2PKHTemplate) IsUnlockingScript(s *script.Script) bool {
	pos := 0

	if op, err := s.ReadOp(&pos); err != nil {
		return false
	} else if op.OpCode != script.Op0 && len(op.Data) == 0 {
		return false
	} else if op, err := s.ReadOp(&pos); err != nil {
		return false
	} else if op.OpCode != script.Op0 && len(op.Data) == 0 {
		return false
	} else if len(*s) > pos+1 {
		return false
	}

	return true
}

func (p P2PKHTemplate) Lock() (*script.Script, error) {
	b := []byte{
		script.OpDUP,
		script.OpHASH160,
		script.OpDATA20,
	}
	b = append(b, p.PKHash...)
	b = append(b, script.OpEQUALVERIFY)
	b = append(b, script.OpCHECKSIG)

	s := script.Script(b)
	return &s, nil
}
