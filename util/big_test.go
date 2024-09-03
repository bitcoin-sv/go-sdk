package util_test

import (
	"math/big"
	"testing"

	"github.com/bitcoin-sv/go-sdk/util"
	"github.com/stretchr/testify/require"
)

func TestBigUmod(t *testing.T) {
	require.Equal(t, int64(2), util.Umod(big.NewInt(-178), big.NewInt(10)).Int64())
}

func TestBigRandomInt(t *testing.T) {
	b := util.NewRandomBigInt(32)
	require.NotNil(t, b)
	require.Len(t, b.Bytes(), 32)
}
