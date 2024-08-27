package primitives

import (
	"bytes"
	e "crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	base58 "github.com/bitcoin-sv/go-sdk/compat/base58"
	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	shamir "github.com/bitcoin-sv/go-sdk/primitives/shamir"
)

var (
	// ErrChecksumMismatch describes an error where decoding failed due
	// to a bad checksum.
	ErrChecksumMismatch = errors.New("checksum mismatch")

	// ErrMalformedPrivateKey describes an error where a WIF-encoded private
	// key cannot be decoded due to being improperly formatted.  This may occur
	// if the byte length is incorrect or an unexpected magic number was
	// encountered.
	ErrMalformedPrivateKey = errors.New("malformed private key")
)

type Network byte

var (
	MainNet Network = 0x80
	TestNet Network = 0xef
)

// compressMagic is the magic byte used to identify a WIF encoding for
// an address created from a compressed serialized public key.
const compressMagic byte = 0x01

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
func NewPrivateKey() (*PrivateKey, error) {
	key, err := e.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(key), nil
}

// PrivateKeyFromHex returns a private key from a hex string.
func PrivateKeyFromHex(privKeyHex string) (*PrivateKey, error) {
	if len(privKeyHex) == 0 {
		return nil, errors.New("private key hex is empty")
	}
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, err
	}
	privKey, _ := PrivateKeyFromBytes(privKeyBytes)
	return privKey, nil
}

// PrivateKeyFromWif returns a private key from a WIF string.
func PrivateKeyFromWif(wif string) (*PrivateKey, error) {
	decoded := base58.Decode(wif)
	decodedLen := len(decoded)
	var compress bool

	// Length of base58 decoded WIF must be 32 bytes + an optional 1 byte
	// (0x01) if compressed, plus 1 byte for netID + 4 bytes of checksum.
	switch decodedLen {
	case 1 + PrivateKeyBytesLen + 1 + 4:
		if decoded[33] != compressMagic {
			return nil, ErrMalformedPrivateKey
		}
		compress = true
	case 1 + PrivateKeyBytesLen + 4:
		compress = false
	default:
		return nil, ErrMalformedPrivateKey
	}

	// Checksum is first four bytes of double SHA256 of the identifier byte
	// and privKey.  Verify this matches the final 4 bytes of the decoded
	// private key.
	var tosum []byte
	if compress {
		tosum = decoded[:1+PrivateKeyBytesLen+1]
	} else {
		tosum = decoded[:1+PrivateKeyBytesLen]
	}
	cksum := crypto.Sha256d(tosum)[:4]
	if !bytes.Equal(cksum, decoded[decodedLen-4:]) {
		return nil, ErrChecksumMismatch
	}

	// netID := decoded[0]
	privKeyBytes := decoded[1 : 1+PrivateKeyBytesLen]
	privKey, _ := PrivateKeyFromBytes(privKeyBytes)
	return privKey, nil
}

func (p *PrivateKey) Wif() string {
	return p.WifPrefix(byte(MainNet))
}

func (p *PrivateKey) WifPrefix(prefix byte) string {
	// Precalculate size.  Maximum number of bytes before base58 encoding
	// is one byte for the network, 32 bytes of private key, possibly one
	// extra byte if the pubkey is to be compressed, and finally four
	// bytes of checksum.

	encodeLen := 1 + PrivateKeyBytesLen + 4
	// For now, we assume compressed = true
	// if compress {
	encodeLen++
	// }

	a := make([]byte, 0, encodeLen)
	a = append(a, prefix)
	// Pad and append bytes manually, instead of using Serialize, to
	// avoid another call to make.
	b := p.D.Bytes()
	a = paddedAppend(PrivateKeyBytesLen, a, b)

	// For now, we assume compressed = true
	// if compress {
	a = append(a, compressMagic)
	// }
	cksum := crypto.Sha256d(a)[:4]
	a = append(a, cksum...)
	return base58.Encode(a)
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

// PrivateKeyBytesLen defines the length in bytes of a serialized private key.
const PrivateKeyBytesLen = 32

// Serialize returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p *PrivateKey) Serialize() []byte {
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
	hmac := crypto.Sha256HMAC(invoiceNumberBin, sharedSecret.encode(true))

	newPrivKey := new(big.Int)
	newPrivKey.Add(p.D, new(big.Int).SetBytes(hmac))
	newPrivKey.Mod(newPrivKey, S256().N)
	privKey, _ := PrivateKeyFromBytes(newPrivKey.Bytes())
	return privKey, nil
}

func (p *PrivateKey) ToPolynomial(threshold int) (*shamir.Polynomial, error) {
	// Check for invalid threshold
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}

	curve := shamir.NewCurve()
	points := make([]*shamir.PointInFiniteField, threshold)

	// Set the first point to (0, key)
	keyValue := p.D
	points[0] = shamir.NewPointInFiniteField(big.NewInt(0), new(big.Int).Set(keyValue))

	// Generate random points for the rest of the polynomial
	for i := 1; i < threshold; i++ {
		x, err := rand.Int(rand.Reader, curve.P)
		if err != nil {
			return nil, err
		}
		y, err := rand.Int(rand.Reader, curve.P)
		if err != nil {
			return nil, err
		}
		points[i] = shamir.NewPointInFiniteField(x, y)
	}

	return shamir.NewPolynomial(points, threshold), nil
}

/**
 * Splits the private key into shares using Shamir's Secret Sharing Scheme.
 *
 * @param threshold The minimum number of shares required to reconstruct the private key.
 * @param totalShares The total number of shares to generate.
 * @param prime The prime number to be used in Shamir's Secret Sharing Scheme.
 * @returns An array of shares.
 *
 * @example
 * const key = PrivateKey.fromRandom()
 * const shares = key.toKeyShares(2, 5)
 */
func (p *PrivateKey) ToKeyShares(threshold int, totalShares int) (keyShares *shamir.KeyShares, error error) {
	// if (typeof threshold !== 'number' || typeof totalShares !== 'number') throw new Error('threshold and totalShares must be numbers')
	// if (threshold < 2) throw new Error('threshold must be at least 2')
	// if (totalShares < 2) throw new Error('totalShares must be at least 2')
	// if (threshold > totalShares) throw new Error('threshold should be less than or equal to totalShares')

	// const poly = Polynomial.fromPrivateKey(this, threshold)

	// const points = []
	// for (let i = 0; i < totalShares; i++) {
	// 	const x = new BigNumber(PrivateKey.fromRandom().toArray())
	// 	const y = poly.valueAt(x)
	// 	points.push(new PointInFiniteField(x, y))
	// }

	// const integrity = (this.toPublicKey().toHash('hex') as string).slice(0, 8)

	// return new KeyShares(points, threshold, integrity)

	// TODO: Port typescript above to go
	if threshold < 2 {
		return nil, errors.New("threshold must be at least 2")
	}
	if totalShares < 2 {
		return nil, errors.New("totalShares must be at least 2")
	}
	if threshold > int(totalShares) {
		return nil, errors.New("threshold should be less than or equal to totalShares")
	}

	poly, err := p.ToPolynomial(threshold)
	if err != nil {
		return nil, err
	}

	var points []*shamir.PointInFiniteField
	for range totalShares {
		pk, err := NewPrivateKey()
		if err != nil {
			return nil, err
		}
		x := new(big.Int)
		x.Set(pk.D)

		y := new(big.Int)
		y.Set(poly.ValueAt(x))
		points = append(points, shamir.NewPointInFiniteField(x, y))
	}

	integrity := hex.EncodeToString(p.PubKey().encode(true))[:16]
	return shamir.NewKeyShares(points, threshold, integrity), nil
}

/**
 * Combines shares to reconstruct the private key.
 *
 * @param shares An array of points (shares) to be used to reconstruct the private key.
 * @param threshold The minimum number of shares required to reconstruct the private key.
 *
 * @returns The reconstructed private key.
 *
 * @example
 * const share1 = '2NWeap6SDBTL5jVnvk9yUxyfLqNrDs2Bw85KNDfLJwRT.4yLtSm327NApsbuP7QXVW3CWDuBRgmS6rRiFkAkTukic'
 * const share2 = '7NbgGA8iAsxg2s6mBLkLFtGKQrnc4aCbooHJJV31cWs4.GUgXtudthawE3Eevc1waT3Atr1Ft7j1XxdUguVo3B7x3'
 * const reconstructedKey = PrivateKey.fromKeyShares({ shares: [share1, share2], threshold: 2, integrity: '23409547' })
 *
 **/
func PrivateKeyFromKeyShares(keyShares *shamir.KeyShares) (*PrivateKey, error) {
	// convert from ts to go
	if keyShares.Threshold < 2 || keyShares.Threshold > 99 {
		return nil, errors.New("threshold should be between 2 and 99")
	}

	if len(keyShares.Points) < keyShares.Threshold {
		return nil, fmt.Errorf("at least %d shares are required to reconstruct the private key", keyShares.Threshold)
	}

	// check to see if two points have the same x value
	for i := 0; i < keyShares.Threshold; i++ {
		for j := i + 1; j < keyShares.Threshold; j++ {
			fmt.Printf("Comparing X values: Point %d (X: %s) and Point %d (X: %s)\n", i, keyShares.Points[i].X.String(), j, keyShares.Points[j].X.String())
			if keyShares.Points[i].X.Cmp(keyShares.Points[j].X) == 0 {
				fmt.Printf("Detected duplicate X value at indices %d and %d with values %s and %s\n", i, j, keyShares.Points[i].X.String(), keyShares.Points[j].X.String())
				return nil, fmt.Errorf("duplicate share detected, each must be unique: %d (%s) and %d (%s)", i, keyShares.Points[i].X, j, keyShares.Points[j].X)
			}
		}
	}

	poly := shamir.NewPolynomial(keyShares.Points, keyShares.Threshold)
	polyBytes := poly.ValueAt(big.NewInt(0)).Bytes()
	privateKey, publicKey := PrivateKeyFromBytes(polyBytes)
	integrityHash := publicKey.encode(true)[:8]
	if keyShares.Integrity != hex.EncodeToString(integrityHash) {
		return nil, fmt.Errorf("integrity hash mismatch %s != %s", keyShares.Integrity, hex.EncodeToString(integrityHash))
	}
	return privateKey, nil
}

//  static fromKeyShares (keyShares: KeyShares): PrivateKey {
//   const { points, threshold, integrity } = keyShares
//   if (threshold < 2 || threshold > 99) throw new Error('threshold should be between 2 and 99')
//   if (points.length < threshold) throw new Error(`At least ${threshold} shares are required to reconstruct the private key`)
//   // check to see if two points have the same x value
//   for (let i = 0; i < threshold; i++) {
//     for (let j = i + 1; j < threshold; j++) {
//       if (points[i].x.eq(points[j].x)) {
//         throw new Error('Duplicate share detected, each must be unique.')
//       }
//     }
//   }
//   const poly = new Polynomial(points, threshold)
//   const privateKey = new PrivateKey(poly.valueAt(new BigNumber(0)).toArray())
//   const integrityHash = privateKey.toPublicKey().toHash('hex').slice(0, 8)
//   if (integrityHash !== integrity) {
//     throw new Error('Integrity hash mismatch')
//   }

//   return privateKey
// }

/**
 * @method toBackupShares
 *
 * Creates a backup of the private key by splitting it into shares.
 *
 *
 * @param threshold The number of shares which will be required to reconstruct the private key.
 * @param totalShares The number of shares to generate for distribution.
 * @returns
 */
func (p *PrivateKey) ToBackupShares(threshold int, shares int) ([]string, error) {
	keyShares, err := p.ToKeyShares(threshold, shares)
	if err != nil {
		return nil, err
	}
	return keyShares.ToBackupFormat()
}

/**
 *
 * @method fromBackupShares
 *
 * Creates a private key from backup shares.
 *
 * @param shares
 * @returns PrivateKey
 */
func PrivateKeyFromBackupShares(shares []string) (*PrivateKey, error) {
	keyShares, err := shamir.NewKeySharesFromBackupFormat(shares)
	if err != nil {
		return nil, err
	}
	return PrivateKeyFromKeyShares(keyShares)
}
