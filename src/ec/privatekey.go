package ec

import (
	e "crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/bitcoin-sv/go-sdk/src/src/crypto"
)

// PrivateKey wraps an ecdsa.PrivateKey as a convenience mainly for signing
// things with the the private key without having to directly import the ecdsa
// package.
type PrivateKey e.PrivateKey

// PrivateKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivateKeyFromBytes(pk []byte) (*PrivateKey, *PublicKey) {
	x, y := S256().ScalarBaseMult(pk)
	priv := &e.PrivateKey{
		PublicKey: e.PublicKey{
			Curve: S256(),
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}
	return (*PrivateKey)(priv), (*PublicKey)(&priv.PublicKey)
}

// NewPrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func NewPrivateKey() (*PrivateKey, error) {
	key, err := e.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(key), nil
}

// PrivateKey is an ecdsa.PrivateKey with additional functions to
func PrivateKeyFromString(privKeyHex string) (*PrivateKey, error) {
	privKeyBytes, _ := hex.DecodeString(privKeyHex)
	privKey, _ := PrivateKeyFromBytes(privKeyBytes)
	return privKey, nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

// Sign generates an ECDSA signature for the provided hash (which should be the result
// of hashing a larger message) using the private key. Produced signature
// is deterministic (same message and same key yield the same signature) and canonical
// in accordance with RFC6979 and BIP0062.
func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	return signRFC6979(p, hash)
}

// PrivateKeyBytesLen defines the length in bytes of a serialised private key.
const PrivateKeyBytesLen = 32

// Serialise returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p *PrivateKey) Serialise() []byte {
	b := make([]byte, 0, PrivateKeyBytesLen)
	return paddedAppend(PrivateKeyBytesLen, b, p.D.Bytes())
}

func (p *PrivateKey) DeriveSharedSecret(key *PublicKey) (*PublicKey, error) {
	if !key.Validate() {
		return nil, fmt.Errorf("public key is not on the curve")
	}
	return key.Mul(p.D), nil
}

// Derives a child key with BRC-42
//
// See BRC-42 spec here: https://github.com/bitcoin-sv/BRCs/blob/master/key-derivation/0042.md
func (p *PrivateKey) DeriveChild(pub *PublicKey, invoiceNumber string) (*PrivateKey, error) {
	invoiceNumberBin := []byte(invoiceNumber)
	sharedSecret, err := p.DeriveSharedSecret(pub)
	if err != nil {
		return nil, err
	}
	pubKeyEncoded := sharedSecret.encode(true)
	hmac := crypto.Sha256HMAC(invoiceNumberBin, pubKeyEncoded)

	newPrivKey := new(big.Int)
	newPrivKey.Add(p.D, new(big.Int).SetBytes(hmac))
	newPrivKey.Mod(newPrivKey, S256().N)
	privKey, _ := PrivateKeyFromBytes(newPrivKey.Bytes())
	return privKey, nil
}
