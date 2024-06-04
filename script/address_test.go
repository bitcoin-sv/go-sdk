package script_test

import (
	"encoding/hex"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/stretchr/testify/assert"
)

const testPublicKeyHash = "00ac6144c4db7b5790f343cf0477a65fb8a02eb7"

func TestNewAddressFromString(t *testing.T) {
	t.Parallel()

	t.Run("mainnet", func(t *testing.T) {
		addressMain := "1E7ucTTWRTahCyViPhxSMor2pj4VGQdFMr"

		addr, err := script.NewAddressFromString(addressMain)
		assert.NoError(t, err)
		assert.NotNil(t, addr)

		assert.Equal(t, "8fe80c75c9560e8b56ed64ea3c26e18d2c52211b", addr.PublicKeyHash.String(), addressMain)
		assert.Equal(t, addressMain, addr.AddressString)
	})

	t.Run("testnet", func(t *testing.T) {
		addressTestnet := "mtdruWYVEV1wz5yL7GvpBj4MgifCB7yhPd"

		addr, err := script.NewAddressFromString(addressTestnet)
		assert.NoError(t, err)
		assert.NotNil(t, addr)

		assert.Equal(t, "8fe80c75c9560e8b56ed64ea3c26e18d2c52211b", addr.PublicKeyHash.String(), addressTestnet)
		assert.Equal(t, addressTestnet, addr.AddressString)
	})

	t.Run("short address", func(t *testing.T) {
		shortAddress := "ADD8E55"
		addr, err := script.NewAddressFromString(shortAddress)
		assert.Error(t, err)
		assert.Nil(t, addr)
		assert.EqualError(t, err, "invalid address length for '"+shortAddress+"'")
	})

	t.Run("unsupported address", func(t *testing.T) {
		unsupportedAddress := "27BvY7rFguYQvEL872Y7Fo77Y3EBApC2EK"
		addr, err := script.NewAddressFromString(unsupportedAddress)
		assert.Error(t, err)
		assert.Nil(t, addr)
		assert.EqualError(t, err, "address not supported "+unsupportedAddress)
	})

}

func TestNewAddressFromPublicKeyString(t *testing.T) {
	t.Parallel()

	t.Run("mainnet", func(t *testing.T) {
		addr, err := script.NewAddressFromPublicKeyString(
			"026cf33373a9f3f6c676b75b543180703df225f7f8edbffedc417718a8ad4e89ce",
			true,
		)
		assert.NoError(t, err)
		assert.NotNil(t, addr)

		assert.Equal(t, testPublicKeyHash, addr.PublicKeyHash.String())
		assert.Equal(t, "114ZWApV4EEU8frr7zygqQcB1V2BodGZuS", addr.AddressString)
	})

	t.Run("testnet", func(t *testing.T) {
		addr, err := script.NewAddressFromPublicKeyString(
			"026cf33373a9f3f6c676b75b543180703df225f7f8edbffedc417718a8ad4e89ce",
			false,
		)
		assert.NoError(t, err)
		assert.NotNil(t, addr)

		assert.Equal(t, testPublicKeyHash, addr.PublicKeyHash.String())
		assert.Equal(t, "mfaWoDuTsFfiunLTqZx4fKpVsUctiDV9jk", addr.AddressString)
	})
}

func TestNewAddressFromPublicKey(t *testing.T) {
	t.Parallel()

	pubKeyBytes, err := hex.DecodeString("026cf33373a9f3f6c676b75b543180703df225f7f8edbffedc417718a8ad4e89ce")
	assert.NoError(t, err)

	var pubKey *ec.PublicKey
	pubKey, err = ec.ParsePubKey(pubKeyBytes)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)

	var addr *script.Address
	addr, err = script.NewAddressFromPublicKey(pubKey, true)
	assert.NoError(t, err)
	assert.NotNil(t, addr)

	assert.Equal(t, testPublicKeyHash, addr.PublicKeyHash.String())
	assert.Equal(t, "114ZWApV4EEU8frr7zygqQcB1V2BodGZuS", addr.AddressString)
}

func TestBase58EncodeMissingChecksum(t *testing.T) {
	t.Parallel()

	input, err := hex.DecodeString("0488b21e000000000000000000362f7a9030543db8751401c387d6a71e870f1895b3a62569d455e8ee5f5f5e5f03036624c6df96984db6b4e625b6707c017eb0e0d137cd13a0c989bfa77a4473fd")
	assert.NoError(t, err)

	assert.Equal(t,
		"xpub661MyMwAqRbcF5ivRisXcZTEoy7d9DfLF6fLqpu5GWMfeUyGHuWJHVp5uexDqXTWoySh8pNx3ELW7qymwPNg3UEYHjwh1tpdm3P9J2j4g32",
		script.Base58EncodeMissingChecksum(input),
	)
}
