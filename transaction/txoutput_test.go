package transaction_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/assert"
)

func TestNewP2PKHOutputFromPubKeyHashStr(t *testing.T) {
	t.Parallel()

	t.Run("valid output", func(t *testing.T) {
		// This is the PKH for address mtdruWYVEV1wz5yL7GvpBj4MgifCB7yhPd
		pkh, _ := hex.DecodeString("8fe80c75c9560e8b56ed64ea3c26e18d2c52211b")
		add, err := script.NewAddressFromPublicKeyHash(pkh, true)
		assert.NoError(t, err)
		s, err := p2pkh.Lock(add)
		assert.NoError(t, err)
		assert.Equal(t,
			"76a9148fe80c75c9560e8b56ed64ea3c26e18d2c52211b88ac",
			hex.EncodeToString(*s),
		)
	})
}

func TestNewHashPuzzleOutput(t *testing.T) {
	t.Parallel()

	t.Run("invalid public key", func(t *testing.T) {
		tx := transaction.NewTransaction()
		err := tx.AddHashPuzzleOutput("", "0", uint64(5000))
		assert.Error(t, err)
	})

	t.Run("missing secret and public key", func(t *testing.T) {
		tx := transaction.NewTransaction()
		err := tx.AddHashPuzzleOutput("", "", uint64(5000))

		assert.NoError(t, err)
		assert.Equal(t,
			"a914b472a266d0bd89c13706a4132ccfb16f7c3b9fcb8876a90088ac",
			tx.Outputs[0].LockingScriptHex(),
		)
	})

	t.Run("valid puzzle output", func(t *testing.T) {
		addr, err := script.NewAddressFromString("myFhJggmsaA2S8Qe6ZQDEcVCwC4wLkvC4e")
		assert.NoError(t, err)
		assert.NotNil(t, addr)

		tx := transaction.NewTransaction()
		err = tx.AddHashPuzzleOutput("secret1", addr.PublicKeyHash.String(), uint64(5000))

		assert.NoError(t, err)
		assert.Equal(t,
			"a914d3f9e3d971764be5838307b175ee4e08ba427b908876a914c28f832c3d539933e0c719297340b34eee0f4c3488ac",
			tx.Outputs[0].LockingScriptHex(),
		)
	})
}

func TestNewOpReturnOutput(t *testing.T) {
	t.Parallel()

	data := "On February 4th, 2020 The Return to Genesis was activated to restore the Satoshi Vision for Bitcoin. " +
		"It is locked in irrevocably by this transaction. Bitcoin can finally be Bitcoin again and the miners can " +
		"continue to write the Chronicle of everything. Thank you and goodnight from team SV."
	dataBytes := []byte(data)

	tx := transaction.NewTransaction()
	err := tx.AddOpReturnOutput(dataBytes)
	assert.NoError(t, err)

	script := tx.Outputs[0].LockingScriptHex()
	dataLength := transaction.VarInt(uint64(len(dataBytes))).Bytes()

	assert.Equal(t, "006a4d2201"+hex.EncodeToString(dataBytes), script)
	assert.Equal(t, "fd2201", fmt.Sprintf("%x", dataLength))
}

func TestNewOpReturnPartsOutput(t *testing.T) {
	t.Parallel()

	dataBytes := [][]byte{[]byte("hi"), []byte("how"), []byte("are"), []byte("you")}
	tx := transaction.NewTransaction()
	err := tx.AddOpReturnPartsOutput(dataBytes)
	assert.NoError(t, err)

	assert.Equal(t, "006a02686903686f770361726503796f75", tx.Outputs[0].LockingScriptHex())
}

func TestTx_TotalOutputSatoshis(t *testing.T) {
	t.Parallel()

	t.Run("greater than zero", func(t *testing.T) {
		tx, err := transaction.NewTransactionFromHex("020000000180f1ada3ad8e861441d9ceab40b68ed98f13695b185cc516226a46697cc01f80010000006b483045022100fa3a0f8fa9fbf09c372b7a318fa6175d022c1d782f7b8bc5949a7c8f59ce3f35022005e0e84c26f26d892b484ff738d803a57626679389c8b302939460dab29a5308412103e46b62eea5db5898fb65f7dc840e8a1dbd8f08a19781a23f1f55914f9bedcd49feffffff02dec537b2000000001976a914ba11bcc46ecf8d88e0828ddbe87997bf759ca85988ac00943577000000001976a91418392a59fc1f76ad6a3c7ffcea20cfcb17bda9eb88ac6e000000")
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, uint64((29.89999582+20.00)*1e8), tx.TotalOutputSatoshis())
	})

	t.Run("zero Outputs", func(t *testing.T) {
		tx := transaction.NewTransaction()
		assert.NotNil(t, tx)
		assert.Equal(t, uint64(0), tx.TotalOutputSatoshis())
	})
}
