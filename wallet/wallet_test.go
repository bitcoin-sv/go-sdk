package wallet_test

import (
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWallet_EncryptDecryptCounterparty(t *testing.T) {
	// Create test data
	sampleData := []byte{3, 1, 4, 1, 5, 9}

	// Generate keys
	userKey, err := ec.NewPrivateKey()
	assert.NoError(t, err)
	counterpartyKey, err := ec.NewPrivateKey()
	assert.NoError(t, err)

	// Create wallets with proper initialization
	userWallet := wallet.NewWallet(userKey)
	counterpartyWallet := wallet.NewWallet(counterpartyKey)

	// Define protocol and key ID
	protocol := wallet.WalletProtocol{
		SecurityLevel: wallet.SecurityLevelEveryAppAndCounterparty,
		Protocol:      "tests",
	}

	// Encrypt message
	encryptArgs := &wallet.WalletEncryptArgs{
		WalletEncryptionArgs: wallet.WalletEncryptionArgs{
			ProtocolID: protocol,
			KeyID:      "4",
			Counterparty: wallet.WalletCounterparty{
				Type:         wallet.CounterpartyTypeOther,
				Counterparty: counterpartyKey.PubKey(),
			},
		},
		Plaintext: sampleData,
	}

	encryptResult, err := userWallet.Encrypt(encryptArgs)
	assert.NoError(t, err)
	assert.NotEqual(t, sampleData, encryptResult.Ciphertext)

	// Decrypt message
	decryptArgs := &wallet.WalletDecryptArgs{
		WalletEncryptionArgs: wallet.WalletEncryptionArgs{
			ProtocolID: protocol,
			KeyID:      "4",
			Counterparty: wallet.WalletCounterparty{
				Type:         wallet.CounterpartyTypeOther,
				Counterparty: userKey.PubKey(),
			},
		},
		Ciphertext: encryptResult.Ciphertext,
	}

	decryptResult, err := counterpartyWallet.Decrypt(decryptArgs)
	assert.NoError(t, err)
	assert.Equal(t, sampleData, decryptResult.Plaintext)

	// Test error cases
	t.Run("wrong protocol", func(t *testing.T) {
		wrongProtocolArgs := *decryptArgs
		wrongProtocolArgs.ProtocolID.Protocol = "wrong"
		_, err := counterpartyWallet.Decrypt(&wrongProtocolArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cipher: message authentication failed")
	})

	t.Run("wrong key ID", func(t *testing.T) {
		wrongKeyArgs := *decryptArgs
		wrongKeyArgs.KeyID = "5"
		_, err := counterpartyWallet.Decrypt(&wrongKeyArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cipher: message authentication failed")
	})

	t.Run("wrong counterparty", func(t *testing.T) {
		wrongCounterpartyArgs := *decryptArgs
		wrongCounterpartyArgs.Counterparty.Counterparty = counterpartyKey.PubKey()
		_, err := counterpartyWallet.Decrypt(&wrongCounterpartyArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cipher: message authentication failed")
	})

	t.Run("invalid protocol name", func(t *testing.T) {
		invalidProtocolArgs := *decryptArgs
		invalidProtocolArgs.ProtocolID.Protocol = "x"
		_, err := counterpartyWallet.Decrypt(&invalidProtocolArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "protocol names must be 5 characters or more")
	})

	t.Run("invalid key ID", func(t *testing.T) {
		invalidKeyArgs := *decryptArgs
		invalidKeyArgs.KeyID = ""
		_, err := counterpartyWallet.Decrypt(&invalidKeyArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key IDs must be 1 character or more")
	})

	t.Run("invalid security level", func(t *testing.T) {
		invalidSecurityArgs := *decryptArgs
		invalidSecurityArgs.ProtocolID.SecurityLevel = -1
		_, err := counterpartyWallet.Decrypt(&invalidSecurityArgs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "protocol security level must be 0, 1, or 2")
	})
}
