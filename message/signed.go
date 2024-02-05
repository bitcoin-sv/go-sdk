package message

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"

	"github.com/bitcoin-sv/go-sdk/ec"
)

// TODO - Implement go version
// import PublicKey from '../primitives/PublicKey.js'
// import PrivateKey from '../primitives/PrivateKey.js'
// import Signature from '../primitives/Signature.js'
// import Curve from '../primitives/Curve.js'
// import Random from '../primitives/Random.js'
// import { toBase64, toArray, Reader, toHex } from '../primitives/utils.js'

// const VERSION = '42423301'

// /**
//  * Signs a message from one party to be verified by another, or for verification by anyone, using the BRC-77 message signing protocol.
//  * @param message The message to sign
//  * @param signer The private key of the message signer
//  * @param [verifier] The public key of the person who can verify the message. If not provided, anyone will be able to verify the message signature.
//  *
//  * @returns The message signature.
//  */
// export const sign = (
//   message: number[],
//   signer: PrivateKey,
//   verifier?: PublicKey
// ): number[] => {
//   const recipientAnyone = typeof verifier !== 'object'
//   if (recipientAnyone) {
//     const curve = new Curve()
//     const anyone = new PrivateKey(1)
//     const anyonePoint = curve.g.mul(anyone)
//     verifier = new PublicKey(
//       anyonePoint.x,
//       anyonePoint.y
//     )
//   }
//   const keyID = Random(32)
//   const keyIDBase64 = toBase64(keyID)
//   const invoiceNumber = `2-message signing-${keyIDBase64}`
//   const signingKey = signer.deriveChild(verifier, invoiceNumber)
//   const signature = signingKey.sign(message).toDER()
//   const senderPublicKey = signer.toPublicKey().encode(true)
//   const version = toArray(VERSION, 'hex')
//   return [
//     ...version,
//     ...senderPublicKey,
//     ...(recipientAnyone ? [0] : verifier.encode(true)),
//     ...keyID,
//     ...signature
//   ]
// }

// /**
//  * Verifies a message using the BRC-77 message signing protocol.
//  * @param message The message to verify.
//  * @param sig The message signature to be verified.
//  * @param [recipient] The private key of the message verifier. This can be omitted if the message is verifiable by anyone.
//  *
//  * @returns True if the message is verified.
//  */
// export const verify = (message: number[], sig: number[], recipient?: PrivateKey): boolean => {
//   const reader = new Reader(sig)
//   const messageVersion = toHex(reader.read(4))
//   if (messageVersion !== VERSION) {
//     throw new Error(
//         `Message version mismatch: Expected ${VERSION}, received ${messageVersion}`
//     )
//   }
//   const signer = PublicKey.fromString(toHex(reader.read(33)))
//   const [verifierFirst] = reader.read(1)
//   if (verifierFirst === 0) {
//     recipient = new PrivateKey(1)
//   } else {
//     const verifierRest = reader.read(32)
//     const verifierDER = toHex([verifierFirst, ...verifierRest])
//     if (typeof recipient !== 'object') {
//       throw new Error(`This signature can only be verified with knowledge of a specific private key. The associated public key is: ${verifierDER}`)
//     }
//     const recipientDER = recipient.toPublicKey().encode(true, 'hex') as string
//     if (verifierDER !== recipientDER) {
//       throw new Error(`The recipient public key is ${recipientDER} but the signature requres the recipient to have public key ${verifierDER}`)
//     }
//   }
//   const keyID = toBase64(reader.read(32))
//   const signatureDER = toHex(reader.read(reader.bin.length - reader.pos))
//   const signature = Signature.fromDER(signatureDER, 'hex')
//   const invoiceNumber = `2-message signing-${keyID}`
//   const signingKey = signer.deriveChild(recipient, invoiceNumber)
//   const verified = signingKey.verify(message, signature)
//   return verified
// }

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
		anyone, err := ec.NewPrivateKey(ec.S256())
		if err != nil {
			return nil, err
		}
		anyonePointX, anyonePointY := ec.S256().ScalarMult(anyone.X, anyone.Y, anyone.Serialise())
		verifier = &ec.PublicKey{X: anyonePointX, Y: anyonePointY}
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
	version := []byte(VERSION)

	sig := append(version, senderPublicKey.SerialiseCompressed()...)
	if recipientAnyone {
		sig = append(sig, 0)
	} else {
		sig = append(sig, verifier.SerialiseCompressed()...)
	}
	sig = append(sig, keyID...)
	sig = append(sig, signature.Serialise()...)
	return sig, nil

}

func Verify(message []byte, sig []byte, recipient *ec.PrivateKey) (bool, error) {
	messageVersion := sig[:4]
	if string(messageVersion) != VERSION {
		return false, nil
	}
	pubKeyBytes := sig[4:37]
	signer, err := ec.ParsePubKey(pubKeyBytes, ec.S256())
	if err != nil {
		return false, err
	}
	verifierFirst := sig[37]
	var recipientPub *ec.PrivateKey
	if verifierFirst == 0 {
		recipientPub, err = ec.NewPrivateKey(ec.S256())
		if err != nil {
			return false, err
		}
	} else {
		verifierRest := sig[38:70]
		verifierDER := append([]byte{verifierFirst}, verifierRest...)
		if recipient == nil {
			return false, nil
		}
		recipientDER := recipient.PubKey().SerialiseCompressed()
		if !bytes.Equal(verifierDER, recipientDER) {
			return false, nil
		}
	}
	keyID := sig[70:102]
	signatureDER := sig[102:]
	signature, err := ec.ParseSignature(signatureDER, ec.S256())
	if err != nil {
		return false, err
	}
	keyIDBase64 := base64.StdEncoding.EncodeToString(keyID)
	invoiceNumber := "2-message signing-" + keyIDBase64
	signingKey, err := signer.DeriveChild(recipientPub, invoiceNumber)
	if err != nil {
		return false, err
	}
	verified := signingKey.Verify(message, signature)
	return verified, nil

}
