package primitives

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBigNumberUmod is a unit test for BigNumber.Umod
func TestBigNumberUmod(t *testing.T) {
	require.Equal(t, 2, NewBigNumberFromInt(-178).Umod(NewBigNumberFromInt(10)).ToNumber())
}
