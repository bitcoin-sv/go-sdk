// Copyright (c) 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package primitives

import (
	"encoding/hex"
	"testing"

	chaincfg "github.com/bitcoin-sv/go-sdk/transaction/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeWIF(t *testing.T) {
	mainWif := "Kxfd8ABTYZHBH3y1jToJ2AUJTMVbsNaqQsrkpo9gnnc1JXfBH8mn"
	testWif := "cP2cb5BJycySSVSH7scRPUyN5ao1XpgXUv1DwDcCHuG1ZGkJQxzH"
	privHex := "2b2afb5ea8d8c623acd6744547628988d86787003bb970afe15d471b727db79c"

	mainPriv, _ := PrivateKeyFromWif(mainWif)
	testPriv, _ := PrivateKeyFromWif(testWif)
	assert.Equal(t, mainWif, mainPriv.Wif())
	assert.Equal(t, testWif, mainPriv.WifForChain(&chaincfg.TestNet))
	assert.Equal(t, privHex, hex.EncodeToString(mainPriv.Serialise()))
	assert.Equal(t, mainPriv.Serialise(), testPriv.Serialise())
}
