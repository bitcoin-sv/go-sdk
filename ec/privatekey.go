package ec

import (
	e "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/bitcoin-sv/go-sdk/crypto"
)

// PrivateKey wraps an ecdsa.PrivateKey as a convenience mainly for signing
// things with the the private key without having to directly import the ecdsa
// package.
type PrivateKey e.PrivateKey

// PrivateKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivateKeyFromBytes(curve elliptic.Curve, pk []byte) (*PrivateKey, *PublicKey) {

	x, y := curve.ScalarBaseMult(pk)

	priv := &e.PrivateKey{
		PublicKey: e.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return (*PrivateKey)(priv), (*PublicKey)(&priv.PublicKey)
}

// NewPrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func NewPrivateKey(curve elliptic.Curve) (*PrivateKey, error) {
	key, err := e.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(key), nil
}

// PubKey returns the PublicKey corresponding to this private key.
func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

// ToECDSA returns the private key as a *ecdsa.PrivateKey.
func (p *PrivateKey) ToECDSA() *PrivateKey {
	return (*PrivateKey)(p)
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
	return paddedAppend(PrivateKeyBytesLen, b, p.ToECDSA().D.Bytes())
}

func (p *PrivateKey) deriveSharedSecret(key *PublicKey) *PrivateKey {
	if !key.Validate() {
		panic("Public key is not on the curve")
	}
	// x or y can be used as the shared secret
	x, _ := SharedSecret(p, key)
	priv, _ := PrivateKeyFromBytes(S256(), x)
	return priv
	// pBigint := new(big.Int).SetBytes(p.Serialise())
	// return key.Mul(pBigint)
}

// Derives a child key with BRC-42
func (p *PrivateKey) DeriveChild(pub *PublicKey, invoiceNumber string) *PrivateKey {
	// return p.deriveSharedSecret(pub)
	sharedSecret := p.deriveSharedSecret(pub)
	invoiceNumberBin := []byte(invoiceNumber)
	pubKeyEncoded, _ := sharedSecret.PubKey().encode(true)
	hmac := crypto.Sha256HMAC(pubKeyEncoded, invoiceNumberBin)
	curve := S256()
	hmacBigint := new(big.Int).SetBytes(hmac)
	privBigint := new(big.Int).SetBytes(p.Serialise())
	privBigint.Add(privBigint, hmacBigint)
	privBigint.Mod(privBigint, curve.Params().N)
	privKey, pubKey := PrivateKeyFromBytes(curve, privBigint.Bytes())
	if !pubKey.Validate() {
		panic("Public key is not on the curve")
	}
	return privKey
}

func SharedSecret(privKeyA *PrivateKey, pubKeyB *PublicKey) ([]byte, []byte) {
	curve := S256()
	x, y := curve.ScalarMult(pubKeyB.X, pubKeyB.Y, privKeyA.D.Bytes())
	return x.Bytes(), y.Bytes()
}
