package template

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	hash "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type P2PKH struct {
	PKHash     []byte
	privateKey *ec.PrivateKey
}

func NewP2PKHFromAddress(address *script.Address) *P2PKH {
	return &P2PKH{
		PKHash: address.PublicKeyHash,
	}
}

func NewP2PKHFromAddressString(addressStr string) (*P2PKH, error) {
	add, err := script.NewAddressFromString(addressStr)
	if err != nil {
		return nil, err
	}
	return NewP2PKHFromAddress(add), nil
}

func NewP2PKHFromPubKey(pubKey []byte) *P2PKH {
	return &P2PKH{
		PKHash: hash.Hash160(pubKey),
	}
}

func NewP2PKHFromPubKeyEC(pubKey *ec.PublicKey) *P2PKH {
	return &P2PKH{
		PKHash: hash.Hash160(pubKey.SerialiseCompressed()),
	}
}

func NewP2PKHFromPrivKey(privKey *ec.PrivateKey) *P2PKH {
	return &P2PKH{
		PKHash:     hash.Hash160(privKey.PubKey().SerialiseCompressed()),
		privateKey: privKey,
	}
}

func (p *P2PKH) IsLockingScript(s *script.Script) bool {
	return s.IsP2PKH()
}

func (p *P2PKH) IsUnlockingScript(s *script.Script) bool {
	pos := 0

	if op, err := s.ReadOp(&pos); err != nil {
		return false
	} else if op.Op != script.Op0 && len(op.Data) == 0 {
		return false
	} else if op, err := s.ReadOp(&pos); err != nil {
		return false
	} else if op.Op != script.Op0 && len(op.Data) == 0 {
		return false
	} else if len(*s) > pos+1 {
		return false
	}

	return true
}

func (p *P2PKH) NewP2PKHFromScript(s *script.Script) *P2PKH {
	p2pkh := &P2PKH{}
	if !p2pkh.IsLockingScript(s) {
		return nil
	}
	p2pkh.PKHash = (*s)[3:23]
	return p2pkh
}

func (p *P2PKH) Lock() (*script.Script, error) {
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

func (p *P2PKH) Unlocker() (transaction.Unlocker, error) {
	return p, nil
}

func (p *P2PKH) Unlock(tx *transaction.Transaction, params *transaction.UnlockParams) (*script.Script, error) {
	if p.privateKey == nil {
		return nil, ErrNoPrivateKey
	}
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}

	if tx.Inputs[params.InputIdx].SourceTransaction == nil {
		return nil, transaction.ErrEmptyPreviousTx
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

	sigBuf := make([]byte, 0)
	sigBuf = append(sigBuf, signature...)
	sigBuf = append(sigBuf, uint8(params.SigHashFlags))

	s := &script.Script{}
	if err = s.AppendPushData(sigBuf); err != nil {
		return nil, err
	} else if err = s.AppendPushData(pubKey); err != nil {
		return nil, err
	}

	return s, nil
}

func (p *P2PKH) EstimateLength(_ *transaction.Transaction, inputIndex uint32) uint32 {
	return 106
}
