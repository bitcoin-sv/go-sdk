package primitives

import "math/big"

type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

type Curve interface {
	Params() *CurveParams
	IsOnCurve(x, y *big.Int) bool
	Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int)
	Double(x1, y1 *big.Int) (*big.Int, *big.Int)
	ScalarMult(x1, y1 *big.Int, k []byte) (*big.Int, *big.Int)
	ScalarBaseMult(k []byte) (*big.Int, *big.Int)
}
