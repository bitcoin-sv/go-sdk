package primitives

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// Curve represents the parameters of the elliptic curve
type Curve struct {
	P *big.Int
}

// func NewCurve() *Curve {
// 	return &Curve{P: big.NewInt(65537)} // 2^16 + 1, a Fermat prime
// }

func NewCurve() *Curve {
	// This is a 256-bit prime number
	p, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	return &Curve{P: p}
}

// PointInFiniteField represents a point in a finite field
type PointInFiniteField struct {
	X, Y *big.Int
}

func NewPointInFiniteField(x, y *big.Int) *PointInFiniteField {
	curve := NewCurve()
	return &PointInFiniteField{
		X: new(big.Int).Mod(x, curve.P),
		Y: new(big.Int).Mod(y, curve.P),
	}
}

func (p *PointInFiniteField) String() string {
	return base64.StdEncoding.EncodeToString(p.X.Bytes()) + "." + base64.StdEncoding.EncodeToString(p.Y.Bytes())
}

func PointFromString(s string) (*PointInFiniteField, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid point string")
	}
	x, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	y, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	return NewPointInFiniteField(new(big.Int).SetBytes(x), new(big.Int).SetBytes(y)), nil
}

type Polynomial struct {
	Points    []*PointInFiniteField
	Threshold int
}

func NewPolynomial(points []*PointInFiniteField, threshold int) *Polynomial {
	if threshold == 0 {
		threshold = len(points)
	}
	return &Polynomial{
		Points:    points,
		Threshold: threshold,
	}
}

func PolynomialFromPrivateKey(privateKey *ec.PrivateKey, threshold int) (*Polynomial, error) {
	// Check for invalid threshold
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}

	curve := NewCurve()
	points := make([]*PointInFiniteField, threshold)

	// Set the first point to (0, key)
	keyValue := privateKey.D
	points[0] = NewPointInFiniteField(big.NewInt(0), new(big.Int).Set(keyValue))

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
		points[i] = NewPointInFiniteField(x, y)
	}

	return NewPolynomial(points, threshold), nil
}

func (p *Polynomial) ValueAt(x *big.Int) *big.Int {
	curve := NewCurve()
	y := big.NewInt(0)

	for i := 0; i < p.Threshold; i++ {
		term := new(big.Int).Set(p.Points[i].Y)
		for j := 0; j < p.Threshold; j++ {
			if i != j {
				numerator := new(big.Int).Sub(x, p.Points[j].X)
				numerator.Mod(numerator, curve.P)

				denominator := new(big.Int).Sub(p.Points[i].X, p.Points[j].X)
				denominator.Mod(denominator, curve.P)

				denominatorInv := new(big.Int).ModInverse(denominator, curve.P)

				fraction := new(big.Int).Mul(numerator, denominatorInv)
				fraction.Mod(fraction, curve.P)

				term.Mul(term, fraction)
				term.Mod(term, curve.P)
			}
		}
		y.Add(y, term)
		y.Mod(y, curve.P)
	}

	return y
}
