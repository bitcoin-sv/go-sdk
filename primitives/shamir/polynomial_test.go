package primitives

import (
	"log"
	"math/big"
	"testing"

	bignumber "github.com/bitcoin-sv/go-sdk/primitives/bignumber"
	"github.com/stretchr/testify/require"
)

// TestPointInFiniteField verifies the creation and string conversion of points
func TestPointInFiniteField(t *testing.T) {
	x := big.NewInt(10)
	y := big.NewInt(20)

	point := NewPointInFiniteField(x, y)

	if point.X.Cmp(x) != 0 || point.Y.Cmp(y) != 0 {
		t.Errorf("Point creation failed. Expected (%v, %v), got (%v, %v)", x, y, point.X, point.Y)
	}

	str := point.String()
	log.Printf("Point string: %s", str)
	reconstructedPoint, err := PointFromString(str)
	if err != nil {
		t.Errorf("Failed to reconstruct point from string: %v", err)
	}

	if reconstructedPoint.X.Cmp(point.X) != 0 || reconstructedPoint.Y.Cmp(point.Y) != 0 {
		t.Errorf("Point reconstruction failed. Expected (%v, %v), got (%v, %v)",
			point.X, point.Y, reconstructedPoint.X, reconstructedPoint.Y)
	}
}

// TestPolynomialValueAt verifies the polynomial evaluation at specific points
func TestPolynomialValueAt(t *testing.T) {
	points := []*PointInFiniteField{
		NewPointInFiniteField(big.NewInt(0), big.NewInt(1)),
		NewPointInFiniteField(big.NewInt(1), big.NewInt(2)),
		NewPointInFiniteField(big.NewInt(2), big.NewInt(3)),
	}
	poly := NewPolynomial(points, 3)

	testCases := []struct {
		x int64
	}{
		{0},
		{1},
		{2},
		{3},
	}

	for _, tc := range testCases {
		x := big.NewInt(tc.x)
		result := poly.ValueAt(x)
		// We're not checking against specific values here, just ensuring it doesn't panic
		t.Logf("Value at x=%d: %v", tc.x, result)
	}
}

// Test BigNumber.Umod with Curve.P
func TestUmod(t *testing.T) {
	require.Equal(t, 100, bignumber.NewBigNumberFromInt(100).Umod(bignumber.NewBigNumber(NewCurve().P)).ToNumber())
	require.Equal(t, 63, bignumber.NewBigNumber(NewCurve().P).Umod(bignumber.NewBigNumberFromInt(100)).ToNumber())
	require.Equal(t, 63, bignumber.NewBigNumber(NewCurve().P).Umod(bignumber.NewBigNumberFromInt(-100)).ToNumber())
}
