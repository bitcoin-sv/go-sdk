package message

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// BRC-78: https://github.com/bitcoin-sv/BRCs/blob/master/peer-to-peer/0078.md
const VERSION = "42423301"

var VERSION_BYTES = []byte{0x42, 0x42, 0x10, 0x33}

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

	priv := ec.NewSymmetricKey(sharedSecret.Compressed()[1:])
	skey := ec.NewSymmetricKey(priv.ToBytes())
	ciphertext, err := skey.Encrypt(message)
	if err != nil {
		return nil, err
	}

	version, err := hex.DecodeString(VERSION)
	if err != nil {
		return nil, err
	}

	senderPublicKey := sender.PubKey().Compressed()
	recipientDER := recipient.Compressed()

	encryptedMessage := append(version, senderPublicKey...)
	encryptedMessage = append(encryptedMessage, recipientDER...)
	encryptedMessage = append(encryptedMessage, keyID[:8]...)
	encryptedMessage = append(encryptedMessage, ciphertext...)
	return encryptedMessage, nil
}

// /**
//   - Decrypts a message from one party to another using the BRC-78 message encryption protocol.
//   - @param message The message to decrypt
//   - @param recipient The private key of the recipient
//     *
//   - @returns The decrypted message
//     */
func Decrypt(message []byte, recipient *ec.PrivateKey) ([]byte, error) {
	messageVersion := message[:4]
	if hex.EncodeToString(messageVersion) != VERSION {
		errorStr := "message version mismatch: Expected %s, received %s"
		return nil, fmt.Errorf(errorStr, VERSION, hex.EncodeToString(messageVersion))
	}
	reader := bytes.NewReader(message)
	_, err := reader.Seek(4, io.SeekStart)
	if err != nil {
		return nil, err
	}
	senderPublicKey := make([]byte, 33)
	_, err = io.ReadFull(reader, senderPublicKey)
	if err != nil {
		return nil, err
	}
	sender, err := ec.ParsePubKey(senderPublicKey)
	if err != nil {
		return nil, err
	}

	expectedRecipientDER := make([]byte, 33)
	_, err = io.ReadFull(reader, expectedRecipientDER)
	if err != nil {
		return nil, err
	}
	actualRecipientDER := recipient.PubKey().Compressed()
	if !bytes.Equal(expectedRecipientDER, actualRecipientDER) {
		errorStr := "the encrypted message expects a recipient public key of %x, but the provided key is %x"
		return nil, fmt.Errorf(errorStr, expectedRecipientDER, actualRecipientDER)
	}
	keyID := make([]byte, 8)
	_, err = io.ReadFull(reader, keyID)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, reader.Len())
	_, err = io.ReadFull(reader, encrypted)
	if err != nil {
		return nil, err
	}
	invoiceNumber := "2-message encryption-" + hex.EncodeToString(keyID)
	signingPub, err := sender.DeriveChild(recipient, invoiceNumber)
	if err != nil {
		return nil, err
	}
	recipientPub, err := recipient.DeriveChild(sender, invoiceNumber)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := signingPub.DeriveSharedSecret(recipientPub)
	if err != nil {
		return nil, err
	}

	priv := ec.NewSymmetricKey(sharedSecret.Compressed()[1:])
	skey := ec.NewSymmetricKey(priv.ToBytes())
	return skey.Decrypt(encrypted)
}
