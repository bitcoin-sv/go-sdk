package message

import (
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/stretchr/testify/require"
	// Using testify for assertions similar to JavaScript's expect
)

func TestSignedMessage(t *testing.T) {
	t.Run("Signs a message for a recipient", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})
		recipientPriv, recipientPub := ec.PrivateKeyFromBytes([]byte{21})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, err := Sign(message, senderPriv, recipientPub)
		require.NoError(t, err)

		verified, err := Verify(message, signature, recipientPriv)
		require.NoError(t, err)
		require.True(t, verified)
	})

	t.Run("Signs a message for anyone", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, err := Sign(message, senderPriv, nil)
		require.NoError(t, err)

		verified, err := Verify(message, signature, nil)
		require.NoError(t, err)
		require.True(t, verified)
	})

	t.Run("Fails to verify a message with a wrong version", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})
		recipientPriv, recipientPub := ec.PrivateKeyFromBytes([]byte{21})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, _ := Sign(message, senderPriv, recipientPub)
		signature[0] = 1 // Alter the signature to simulate version mismatch

		_, err := Verify(message, signature, recipientPriv)
		require.Error(t, err)
		require.Equal(t, "message version mismatch: Expected 42421033, received 01421033", err.Error())
	})

}
