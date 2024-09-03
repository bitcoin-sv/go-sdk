package primitives

import (
	"crypto/rand"
	"math/big"
)

type BigNumber struct {
	value *big.Int
}

func NewBigNumberFromInt64(value int64) *BigNumber {
	return &BigNumber{
		value: big.NewInt(value),
	}
}

func NewBigNumber(value *big.Int) *BigNumber {
	return &BigNumber{
		value: big.NewInt(0).Set(value),
	}
}

func NewBigNumberFromInt(value int64) *BigNumber {
	return NewBigNumber(big.NewInt(value))
}

// Umod returns the unsigned modulus of x and y.
// It ensures the result is always non-negative.
func (x *BigNumber) Umod(y *BigNumber) *BigNumber {
	// Ensure divisor y is not zero
	if y.value.Sign() == 0 {
		panic("division by zero")
	}

	// If the dividend x is zero, the modulus is zero
	if x.value.Sign() == 0 {
		return NewBigNumber(big.NewInt(0))
	}

	mod := new(big.Int)

	// Handle cases where x is negative and y is positive
	if x.value.Sign() < 0 && y.value.Sign() > 0 {
		mod.Neg(x.value)
		mod.Mod(mod, y.value)
		mod.Neg(mod)
		if mod.Sign() < 0 {
			mod.Add(mod, y.value)
		}
	} else if x.value.Sign() > 0 && y.value.Sign() < 0 {
		// Handle cases where x is positive and y is negative
		mod.Mod(x.value, new(big.Int).Neg(y.value))
	} else if x.value.Sign() < 0 && y.value.Sign() < 0 {
		// Handle cases where both x and y are negative
		mod.Neg(x.value)
		mod.Mod(mod, new(big.Int).Neg(y.value))
		mod.Neg(mod)
		if mod.Sign() < 0 {
			mod.Sub(mod, y.value)
		}
	} else {
		// Both numbers are positive
		mod.Mod(x.value, y.value)
	}

	// Ensure the result is non-negative
	if mod.Sign() < 0 {
		mod.Add(mod, new(big.Int).Abs(y.value))
	}

	return NewBigNumber(mod)
}

func (a *BigNumber) ToNumber() int {
	return int(a.value.Int64())
}

func (a *BigNumber) ToBigInt() *big.Int {
	return a.value
}

func Random(byteLen int) *BigNumber {
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return &BigNumber{
		value: new(big.Int).SetBytes(b),
	}
}
