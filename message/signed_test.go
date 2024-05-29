package message

import (
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert" // Using testify for assertions similar to JavaScript's expect
)

func TestSignedMessage(t *testing.T) {
	t.Run("Signs a message for a recipient", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})
		recipientPriv, recipientPub := ec.PrivateKeyFromBytes([]byte{21})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, err := Sign(message, senderPriv, recipientPub)
		assert.Nil(t, err)

		verified, err := Verify(message, signature, recipientPriv)
		assert.Nil(t, err)
		assert.True(t, verified)
	})

	t.Run("Signs a message for anyone", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, err := Sign(message, senderPriv, nil)
		assert.Nil(t, err)

		verified, err := Verify(message, signature, nil)
		assert.Nil(t, err)
		assert.True(t, verified)
	})

	t.Run("Fails to verify a message with a wrong version", func(t *testing.T) {
		senderPriv, _ := ec.PrivateKeyFromBytes([]byte{15})
		recipientPriv, recipientPub := ec.PrivateKeyFromBytes([]byte{21})

		message := []byte{1, 2, 4, 8, 16, 32}
		signature, _ := Sign(message, senderPriv, recipientPub)
		signature[0] = 1 // Alter the signature to simulate version mismatch

		_, err := Verify(message, signature, recipientPriv)
		assert.NotNil(t, err)
		assert.Equal(t, "message version mismatch: Expected 10334242, received 01334242", err.Error())
	})

}
