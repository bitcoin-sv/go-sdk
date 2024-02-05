package message

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/aesgcm"
	"github.com/bitcoin-sv/go-sdk/ec"
)

// BRC-78: https://github.com/bitcoin-sv/BRCs/blob/master/peer-to-peer/0078.md
const VERSION = "42421033"

// Encrypt encrypts a message using the sender's private key and the recipient's public key.
func Encrypt(message []byte, sender *ec.PrivateKey, recipient *ec.PublicKey) ([]byte, error) {
	var keyID [8]byte
	if _, err := rand.Read(keyID[:]); err != nil {
		return nil, err
	}
	keyIDBase64 := hex.EncodeToString(keyID[:])
	invoiceNumber := "2-message encryption-" + keyIDBase64
	signingPriv, err := sender.DeriveChild(recipient, invoiceNumber)
	if err != nil {
		return nil, err
	}
	recipientPub, err := recipient.DeriveChild(sender, invoiceNumber)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := signingPriv.DeriveSharedSecret(recipientPub)
	if err != nil {
		return nil, err
	}
	// FIXME - nonce is not being set correctly
	nonce := make([]byte, 12)
	cyphertext, _, err := aesgcm.EncryptGCM(message, nonce, sharedSecret.SerialiseCompressed(), keyID[:])
	if err != nil {
		return nil, err
	}

	// derive a shared secret
	// sharedSecret, err := signingPriv.ECDH(recipientPub)

	// const symmetricKey = new SymmetricKey(sharedSecret.encode(true).slice(1))
	// const encrypted = symmetricKey.encrypt(message) as number[]
	// const senderPublicKey = sender.toPublicKey().encode(true)
	// const version = toArray(VERSION, 'hex')
	// return [
	//   ...version,
	//   ...senderPublicKey,
	//   ...recipient.encode(true),
	//   ...keyID,
	//   ...encrypted
	// ]

	version, err := hex.DecodeString(VERSION)
	if err != nil {
		return nil, err
	}

	senderPublicKey := sender.PubKey().SerialiseCompressed()
	recipientDER := recipient.SerialiseCompressed()

	encryptedMessage := append(version, senderPublicKey...)
	encryptedMessage = append(encryptedMessage, recipientDER...)
	encryptedMessage = append(encryptedMessage, keyID[:8]...)
	encryptedMessage = append(encryptedMessage, cyphertext...)
	return encryptedMessage, nil
}

// /**
//   - Decrypts a message from one party to another using the BRC-78 message encryption protocol.
//   - @param message The message to decrypt
//   - @param sender The private key of the recipient
//     *
//   - @returns The decrypted message
//     */
func Decrypt(message []byte, recipient *ec.PrivateKey) ([]byte, error) {
	messageVersion := message[:4]
	if string(messageVersion) != VERSION {
		return nil, fmt.Errorf("Message version mismatch: Expected %s, received %s", VERSION, messageVersion)
	}
	reader := bytes.NewReader(message)
	reader.Seek(4, 0)
	senderPublicKey := make([]byte, 33)
	_, err := reader.Read(senderPublicKey)
	if err != nil {
		return nil, err
	}
	sender, err := ec.ParsePubKey(senderPublicKey, ec.S256())
	if err != nil {
		return nil, err
	}

	expectedRecipientDER := make([]byte, 33)
	_, err = reader.Read(expectedRecipientDER)
	if err != nil {
		return nil, err
	}
	actualRecipientDER := recipient.PubKey().SerialiseCompressed()
	if !bytes.Equal(expectedRecipientDER, actualRecipientDER) {
		return nil, fmt.Errorf("the encrypted message expects a recipient public key of %x, but the provided key is %x", expectedRecipientDER, actualRecipientDER)
	}
	keyID := make([]byte, 32)
	_, err = reader.Read(keyID)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, reader.Len())
	_, err = reader.Read(encrypted)
	if err != nil {
		return nil, err
	}
	invoiceNumber := "2-message encryption-" + hex.EncodeToString(keyID)
	signingPriv, err := sender.DeriveChild(recipient, invoiceNumber)
	if err != nil {
		return nil, err
	}
	// TODO: This is weird from ts library - its called recipientPub but its a private key
	recipientPub, err := recipient.DeriveChild(recipient.PubKey(), invoiceNumber)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := signingPriv.DeriveSharedSecret(recipientPub)
	if err != nil {
		return nil, err
	}
	return aesgcm.DecryptGCM(encrypted, sharedSecret.SerialiseCompressed(), keyID, nil, nil)
}
