package message

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

type SignedMessage struct {
	Version            []byte
	SenderPublicKey    *ec.PublicKey
	RecipientPublicKey *ec.PublicKey
	KeyID              []byte
	Signature          *ec.Signature
}

func Sign(message []byte, signer *ec.PrivateKey, verifier *ec.PublicKey) ([]byte, error) {
	recipientAnyone := verifier == nil
	if recipientAnyone {
		_, verifier = ec.PrivateKeyFromBytes([]byte{1})
	}

	keyID := make([]byte, 32)
	_, err := rand.Read(keyID)
	if err != nil {
		return nil, err
	}
	keyIDBase64 := base64.StdEncoding.EncodeToString(keyID)
	invoiceNumber := "2-message signing-" + keyIDBase64
	signingPriv, err := signer.DeriveChild(verifier, invoiceNumber)
	if err != nil {
		return nil, err
	}
	signature, err := signingPriv.Sign(message)
	if err != nil {
		return nil, err
	}
	senderPublicKey := signer.PubKey()

	sig := append(VERSION_BYTES, senderPublicKey.SerializeCompressed()...)
	if recipientAnyone {
		sig = append(sig, 0)
	} else {
		sig = append(sig, verifier.SerializeCompressed()...)
	}
	sig = append(sig, keyID...)
	signatureDER, err := signature.ToDER()
	if err != nil {
		return nil, err
	}
	sig = append(sig, signatureDER...)
	return sig, nil
}

func Verify(message []byte, sig []byte, recipient *ec.PrivateKey) (bool, error) {
	counter := 4
	messageVersion := sig[:counter]
	if !bytes.Equal(messageVersion, VERSION_BYTES) {
		return false, fmt.Errorf("message version mismatch: Expected %x, received %x", VERSION_BYTES, messageVersion)
	}
	pubKeyBytes := sig[counter : counter+33]
	counter += 33
	signer, err := ec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false, err
	}
	verifierFirst := sig[counter]
	if verifierFirst == 0 {
		recipient, _ = ec.PrivateKeyFromBytes([]byte{1})
		counter++
	} else {
		counter++
		verifierRest := sig[counter : counter+32]
		counter += 32
		verifierDER := append([]byte{verifierFirst}, verifierRest...)
		if recipient == nil {
			return false, nil
		}
		recipientDER := recipient.PubKey().SerializeCompressed()
		if !bytes.Equal(verifierDER, recipientDER) {
			err = fmt.Errorf("the recipient public key is %x but the signature requres the recipient to have public key %x", recipientDER, verifierDER)
			return false, err
		}
	}
	keyID := sig[counter : counter+32]
	counter += 32
	signatureDER := sig[counter:]
	signature, err := ec.FromDER(signatureDER)
	if err != nil {
		return false, err
	}
	keyIDBase64 := base64.StdEncoding.EncodeToString(keyID)
	invoiceNumber := "2-message signing-" + keyIDBase64
	signingKey, err := signer.DeriveChild(recipient, invoiceNumber)
	if err != nil {
		return false, err
	}
	verified := signature.Verify(message, signingKey)
	return verified, nil

}
