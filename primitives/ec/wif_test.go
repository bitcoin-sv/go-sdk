// Copyright (c) 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package primitives

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeWIF(t *testing.T) {
	mainWif := "Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn"
	testWif := "cP2cb5BJycySSVSH7scRPUyN5ao1XpgXUv1DwDcCHuG1ZGkJQxzH"
	privHex := "2b2afb5ea8d8c623acd6744547628988d86787003bb970afe15d471b727db79c"

	mainPriv, _ := PrivateKeyFromWif(mainWif)
	testPriv, _ := PrivateKeyFromWif(testWif)
	require.Equal(t, mainWif, mainPriv.Wif())
	require.Equal(t, testWif, mainPriv.WifPrefix(byte(TestNet)))
	require.Equal(t, privHex, hex.EncodeToString(mainPriv.Serialise()))
	require.Equal(t, mainPriv.Serialise(), testPriv.Serialise())
}
