package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/aesgcm"
	"github.com/bitcoin-sv/go-sdk/crypto"
	"github.com/bitcoin-sv/go-sdk/ec"
)

// BRC-78: https://github.com/bitcoin-sv/BRCs/blob/master/peer-to-peer/0078.md

const VERSION = "42421033"

type Message struct {
	version         []byte
	senderPublicKey []byte
	recipient       []byte
	keyID           []byte
	encrypted       []byte
}

// Encrypt encrypts a message using the sender's private key and the recipient's public key.
func Encrypt(message []byte, sender *ec.PrivateKey, recipient *ec.PublicKey) (*Message, error) {
	// secure random 32 bytes
	var keyID [8]byte
	if _, err := rand.Read(keyID[:]); err != nil {
		return nil, err
	}
	// golang - convert to base64
	keyIDBase64 := hex.EncodeToString(keyID[:])
	// TODO: Look into this - TS library is using this pattern. Might be definied in the
	invoiceNumber := "2-message encryption-" + keyIDBase64

	// Derive a signing key
	signingPriv, err := sender.DeriveChild(recipient, invoiceNumber)
	if err != nil {
		return nil, err
	}
	//	const recipientPub = recipient.deriveChild(sender, invoiceNumber)
	//
	// convert from typescript above
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

	return &Message{
		version:         version,
		senderPublicKey: sender.PubKey().SerialiseCompressed(),
		recipient:       recipient.SerialiseCompressed(),
		keyID:           keyID[:],
		encrypted:       cyphertext,
	}, nil

}

// FIXME: This is just sample code - need to implement messaging standard
// deriveMessagePoint takes a message, hashes it, and maps it to a point on the elliptic curve.
func deriveMessagePoint(message []byte) (ec.Point, string, error) {
	// Hash the message using SHA-256
	hash := crypto.Sha256(message)
	cId := hex.EncodeToString(hash[:])

	// Convert hash to a big integer
	// mBn := new(big.Int).SetBytes(hash[:])

	// Get the generator point of the elliptic curve
	curve := ec.S256()
	Gx, Gy := curve.Params().Gx, curve.Params().Gy

	// TODO: Dont multiple by a random hash
	// Multiply the hash big integer with the generator point
	// Mx, My := curve.ScalarMult(Gx, Gy, mBn.Bytes())

	return ec.Point{X: Gx, Y: Gy}, cId, nil
}

func main() {
	// Example usage
	message := []byte("Hello, World!")
	point, cId, err := deriveMessagePoint(message)
	if err != nil {
		panic(err)
	}

	// Output the results
	fmt.Printf("Message Point: (%s, %s)\n", point.X.String(), point.Y.String())
	fmt.Printf("cId: %s\n", cId)
}
