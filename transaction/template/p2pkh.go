package template

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	hash "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type P2PKHTemplate struct {
	PKHash     []byte
	privateKey *ec.PrivateKey
}

func NewP2PKHTemplateFromAddress(address *script.Address) *P2PKHTemplate {
	return &P2PKHTemplate{
		PKHash: address.PublicKeyHash,
	}
}

func NewP2PKHTemplateFromPubKey(pubKey []byte) *P2PKHTemplate {
	return &P2PKHTemplate{
		PKHash: hash.Hash160(pubKey),
	}
}

func NewP2PKHTemplateFromPubKeyEC(pubKey *ec.PublicKey) *P2PKHTemplate {
	return &P2PKHTemplate{
		PKHash: hash.Hash160(pubKey.SerialiseCompressed()),
	}
}

func NewP2PKHTemplateFromPrivKey(privKey *ec.PrivateKey) *P2PKHTemplate {
	return &P2PKHTemplate{
		PKHash:     hash.Hash160(privKey.PubKey().SerialiseCompressed()),
		privateKey: privKey,
	}
}

func (p *P2PKHTemplate) IsLockingScript(s *script.Script) bool {
	return s.IsP2PKH()
}

func (p *P2PKHTemplate) IsUnlockingScript(s *script.Script) bool {
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

func (p *P2PKHTemplate) Lock() (*script.Script, error) {
	if len(p.PKHash) != 20 {
		return nil, ErrBadPublicKeyHash
	}
	b := make([]byte, 0, 25)
	b = append(b, script.OpDUP, script.OpHASH160, script.OpDATA20)
	b = append(b, p.PKHash...)
	b = append(b, script.OpEQUALVERIFY, script.OpCHECKSIG)
	s := script.Script(b)
	return &s, nil
}

func (p *P2PKHTemplate) Sign(tx *transaction.Transaction, params transaction.UnlockParams) (*script.Script, error) {
	if p.privateKey == nil {
		return nil, ErrNoPrivateKey
	}
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}

	if tx.Inputs[params.InputIdx].PreviousTxScript == nil {
		return nil, transaction.ErrEmptyPreviousTxScript
	}

	sh, err := tx.CalcInputSignatureHash(params.InputIdx, params.SigHashFlags)
	if err != nil {
		return nil, err
	}

	sig, err := p.privateKey.Sign(sh)
	if err != nil {
		return nil, err
	}

	pubKey := p.privateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	uscript, err := script.NewP2PKHUnlockingScript(pubKey, signature, params.SigHashFlags)
	if err != nil {
		return nil, err
	}

	return uscript, nil
}

func (p *P2PKHTemplate) EstimateSize(_ *transaction.Transaction, inputIndex uint32) int {
	return 106
}
