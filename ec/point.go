package ec

import (
	"math/big"
)

type Point struct {
	Curve *Curve
	X, Y  *big.Int
}
