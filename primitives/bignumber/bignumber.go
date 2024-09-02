package primitives

import "math/big"

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
		value: value,
	}
}

func NewBigNumberFromInt(value int64) *BigNumber {
	return NewBigNumber(big.NewInt(value))
}

// umod returns the unsigned modulus of a and b.
// It ensures the result is always non-negative.
func (a *BigNumber) Umod(b *BigNumber) *BigNumber {
	div := new(big.Int)
	mod := new(big.Int)
	div.DivMod(a.value, b.value, mod)
	return &BigNumber{
		value: mod,
	}
}

func (a *BigNumber) ToNumber() int {
	return int(a.value.Int64())
}

func (a *BigNumber) ToBigInt() *big.Int {
	return a.value
}
