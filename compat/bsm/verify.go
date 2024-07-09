package bsm

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
)

// PubKeyFromSignature gets a publickey for a signature and tells you whether is was compressed
func PubKeyFromSignature(sig, data string) (pubKey *ec.PublicKey, wasCompressed bool, err error) {

	var decodedSig []byte
	if decodedSig, err = base64.StdEncoding.DecodeString(sig); err != nil {
		return nil, false, err
	}

	// Validate the signature - this just shows that it was valid at all
	// we will compare it with the key next
	var buf bytes.Buffer
	// if err = wire.WriteVarString(&buf, 0, hBSV); err != nil {
	// 	return nil, false, err
	// }
	// if err = wire.WriteVarString(&buf, 0, data); err != nil {
	// 	return nil, false, err
	// }

	varInt := transaction.VarInt(len(hBSV))
	buf.Write(varInt.Bytes())
	// append the hBsv to buff
	buf.WriteString(hBSV)

	varInt = transaction.VarInt(len(data))
	buf.Write(varInt.Bytes())
	// append the data to buff
	buf.WriteString(data)

	// Create the hash
	expectedMessageHash := crypto.Sha256d(buf.Bytes())
	return ec.RecoverCompact(decodedSig, expectedMessageHash)
}

// VerifyMessage verifies a string and address against the provided
// signature and assumes Bitcoin Signed Message encoding.
// The key referenced by the signature must relate to the address provided.
// Do not provide an address from an uncompressed key along with
// a signature from a compressed key
//
// Error will occur if verify fails or verification is not successful (no bool)
// Spec: https://docs.moneybutton.com/docs/bsv-message.html
func VerifyMessage(address, sig, data string) error {

	// Reconstruct the pubkey
	publicKey, wasCompressed, err := PubKeyFromSignature(sig, data)
	if err != nil {
		return err
	}

	// Get the address
	var bscriptAddress *script.Address
	if bscriptAddress, err = script.NewAddressFromPublicKey(publicKey, wasCompressed); err != nil {
		return err
	}

	// Return nil if addresses match.
	if bscriptAddress.AddressString == address {
		return nil
	}
	return fmt.Errorf(
		"address (%s) not found - compressed: %t\n%s was found instead",
		address,
		wasCompressed,
		bscriptAddress.AddressString,
	)
}

// VerifyMessageDER will take a message string, a public key string and a signature string
// (in strict DER format) and verify that the message was signed by the public key.
//
// Copyright (c) 2019 Bitcoin Association
// License: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/LICENSE
//
// Source: https://github.com/bitcoin-sv/merchantapi-reference/blob/master/handler/global.go
func VerifyMessageDER(hash [32]byte, pubKey string, signature string) (verified bool, err error) {

	// Decode the signature string
	var sigBytes []byte
	if sigBytes, err = hex.DecodeString(signature); err != nil {
		return
	}

	// Parse the signature
	var sig *ec.Signature
	if sig, err = ec.ParseDERSignature(sigBytes); err != nil {
		return
	}

	// Decode the pubKey
	var pubKeyBytes []byte
	if pubKeyBytes, err = hex.DecodeString(pubKey); err != nil {
		return
	}

	// Parse the pubKey
	var rawPubKey *ec.PublicKey
	if rawPubKey, err = ec.ParsePubKey(pubKeyBytes); err != nil {
		return
	}

	// Verify the signature against the pubKey
	verified = sig.Verify(hash[:], rawPubKey)
	return
}
