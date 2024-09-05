package compat

import (
	"crypto/hmac"
	"log"
	"math/big"
	"reflect"

	"encoding/base64"
	"errors"

	ecies "github.com/bitcoin-sv/go-sdk/primitives/aescbc"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	c "github.com/bitcoin-sv/go-sdk/primitives/hash"
)

//
// ECIES encryption/decryption methods; AES-128-CBC with PKCS7 is used as the cipher; hmac-sha256 is used as the mac
//

func EncryptSingle(message string, privateKey *ec.PrivateKey) (string, error) {
	messageBytes := []byte(message)
	return ElectrumEncrypt(messageBytes, privateKey.PubKey(), privateKey, false)
}

func DecryptSingle(encryptedData string, privateKey *ec.PrivateKey) (string, error) {
	plainBytes, err := ElectrumDecrypt(encryptedData, privateKey, nil)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

func EncryptShared(message string, toPublicKey *ec.PublicKey, fromPrivateKey *ec.PrivateKey) (string, error) {
	messageBytes := []byte(message)
	return ElectrumEncrypt(messageBytes, toPublicKey, nil, false)
}

func DecryptShared(encryptedData string, toPrivateKey *ec.PrivateKey, fromPublicKey *ec.PublicKey) (string, error) {
	plainBytes, err := ElectrumDecrypt(encryptedData, toPrivateKey, fromPublicKey)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

func ElectrumEncrypt(message []byte, toPublicKey *ec.PublicKey, fromPrivateKey *ec.PrivateKey, noKey bool) (string, error) {
	// Generate an ephemeral EC private key if fromPrivateKey is nil
	var ephemeralPrivateKey *ec.PrivateKey
	if fromPrivateKey == nil {
		ephemeralPrivateKey, _ = ec.NewPrivateKey()
	} else {
		ephemeralPrivateKey = fromPrivateKey
	}

	// Derive ECDH key
	x, y := toPublicKey.Curve.ScalarMult(toPublicKey.X, toPublicKey.Y, ephemeralPrivateKey.D.Bytes())
	ecdhKey := (&ec.PublicKey{X: x, Y: y}).SerializeCompressed()

	// SHA512(ECDH_KEY)
	key := c.Sha512(ecdhKey)
	iv, keyE, keyM := key[0:16], key[16:32], key[32:]

	// AES encryption
	cipherText, err := ecies.AESCBCEncrypt(message, keyE, iv, false)
	if err != nil {
		return "", err
	}

	ephemeralPublicKey := ephemeralPrivateKey.PubKey()
	var encrypted []byte
	if noKey {
		// encrypted = magic_bytes(4 bytes) + cipher
		encrypted = append([]byte("BIE1"), cipherText...)
	} else {
		// encrypted = magic_bytes(4 bytes) + ephemeral_public_key(33 bytes) + cipher
		encrypted = append(append([]byte("BIE1"), ephemeralPublicKey.SerializeCompressed()...), cipherText...)
	}

	mac := c.Sha256HMAC(encrypted, keyM)

	return base64.StdEncoding.EncodeToString(append(encrypted, mac...)), nil
}
func ElectrumDecrypt(encryptedData string, toPrivateKey *ec.PrivateKey, fromPublicKey *ec.PublicKey) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < 52 { // Minimum length: 4 (magic) + 16 (min cipher) + 32 (mac)
		return nil, errors.New("invalid encrypted text: length")
	}
	magic := encrypted[:4]
	if string(magic) != "BIE1" {
		return nil, errors.New("invalid cipher text: invalid magic bytes")
	}

	var sharedSecret []byte
	var cipherText []byte
	var ephemeralPublicKey *ec.PublicKey

	if fromPublicKey != nil {
		// Use counterparty public key to derive shared secret
		x, y := toPrivateKey.Curve.ScalarMult(fromPublicKey.X, fromPublicKey.Y, toPrivateKey.D.Bytes())
		sharedSecret = (&ec.PublicKey{X: x, Y: y}).SerializeCompressed()
		if len(encrypted) > 69 { // 4 (magic) + 33 (pubkey) + 32 (mac)
			cipherText = encrypted[37 : len(encrypted)-32]
		} else {
			cipherText = encrypted[4 : len(encrypted)-32]
		}
	} else {
		// Use ephemeral public key to derive shared secret
		ephemeralPublicKey, err = ec.ParsePubKey(encrypted[4:37])
		if err != nil {
			return nil, err
		}
		x, y := ephemeralPublicKey.Curve.ScalarMult(ephemeralPublicKey.X, ephemeralPublicKey.Y, toPrivateKey.D.Bytes())
		sharedSecret = (&ec.PublicKey{X: x, Y: y}).SerializeCompressed()
		cipherText = encrypted[37 : len(encrypted)-32]
	}

	// Derive key_e, iv and key_m
	key := c.Sha512(sharedSecret)
	iv, keyE, keyM := key[0:16], key[16:32], key[32:]

	// Verify mac
	mac := encrypted[len(encrypted)-32:]
	macRecalculated := c.Sha256HMAC(encrypted[:len(encrypted)-32], keyM)
	if !reflect.DeepEqual(mac, macRecalculated) {
		return nil, errors.New("incorrect password")
	}

	// AES decryption
	plain, err := ecies.AESCBCDecrypt(cipherText, keyE, iv)
	if err != nil {
		return nil, err
	}
	log.Printf("IV: %x, plain text: %x", iv, plain)
	return plain, nil
}

// func BitcoreEncrypt(message []byte, publicKey *ec.PublicKey) (string, error) {
// 	// Generate an ephemeral EC private key
// 	ephemeralPrivateKey, err := ec.NewPrivateKey()
// 	if err != nil {
// 		return "", err
// 	}

// 	// Derive shared secret
// 	x, _ := publicKey.Curve.ScalarMult(publicKey.X, publicKey.Y, ephemeralPrivateKey.D.Bytes())
// 	sharedSecret := x.Bytes()

// 	// Key derivation
// 	keyMaterial := c.Sha512(sharedSecret)

// 	keyE := keyMaterial[:32] // AES-256 key
// 	keyM := keyMaterial[32:] // HMAC key
// 	iv := make([]byte, 16)   // IV for AES (all zeros in Bitcore)

// 	// Encrypt the message
// 	cipherText, err := aescbc.AESEncryptWithIV(message, keyE, iv)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Prepare the output
// 	ephemeralPublicKey := ephemeralPrivateKey.PubKey().SerializeCompressed()
// 	encryptedData := append(ephemeralPublicKey, cipherText...)

// 	// Calculate HMAC
// 	hmacSum := c.Sha256HMAC(encryptedData, keyM)

// 	// Combine all parts
// 	result := append(encryptedData, hmacSum...)

// 	return base64.StdEncoding.EncodeToString(result), nil
// }

func BitcoreEncrypt(message []byte, toPublicKey *ec.PublicKey, fromPrivateKey *ec.PrivateKey, iv []byte) (string, error) {

	// JS Implementation
	// if (!fromPrivateKey) {
	// 	fromPrivateKey = PrivateKey.fromRandom()
	// }
	// const r = fromPrivateKey
	// const RPublicKey = fromPrivateKey.toPublicKey()
	// const RBuf = RPublicKey.encode(true) as number[]
	// const KB = toPublicKey
	// const P = KB.mul(r)
	// const S = P.getX()
	// const Sbuf = S.toArray('be', 32)
	// const kEkM = Hash.sha512(Sbuf)
	// const kE = kEkM.slice(0, 32)
	// const kM = kEkM.slice(32, 64)
	// const c = AESCBC.encrypt(messageBuf, kE, ivBuf)
	// const d = Hash.sha256hmac(kM, [...c])
	// const encBuf = [...RBuf, ...c, ...d]
	// return encBuf
	// If IV is not provided, generate a random one

	if iv == nil {
		iv = make([]byte, 16)
	}

	// If fromPrivateKey is not provided, generate a random one
	if fromPrivateKey == nil {
		var err error
		fromPrivateKey, err = ec.NewPrivateKey()
		if err != nil {
			return "", err
		}
	}

	RBuf := fromPrivateKey.PubKey().ToDERBytes()
	P := toPublicKey.Mul(fromPrivateKey.D)

	Sbuf := P.X.Bytes()
	kEkM := c.Sha512(Sbuf)
	kE := kEkM[:32]
	kM := kEkM[32:]
	cc, err := ecies.AESCBCEncrypt(message, kE, iv, true)
	if err != nil {
		return "", err
	}
	d := c.Sha256HMAC(cc, kM)
	encBuf := append(RBuf, cc...)
	encBuf = append(encBuf, d...)

	log.Printf("encBuf: %x", encBuf)
	result := base64.StdEncoding.EncodeToString(encBuf)
	return result, nil
}

func BitcoreDecrypt(encryptedMessage string, toPrivatKey *ec.PrivateKey) ([]byte, error) {
	// const kB = toPrivateKey
	//   const fromPublicKey = PublicKey.fromString(toHex(encBuf.slice(0, 33)))
	//   const R = fromPublicKey
	//   const P = R.mul(kB)
	//   if (P.eq(new Point(0, 0))) {
	//     throw new Error('P equals 0')
	//   }
	//   const S = P.getX()
	//   const Sbuf = S.toArray('be', 32)
	//   const kEkM = Hash.sha512(Sbuf)
	//   const kE = kEkM.slice(0, 32)
	//   const kM = kEkM.slice(32, 64)
	//   const c = encBuf.slice(33, encBuf.length - 32)
	//   const d = encBuf.slice(encBuf.length - 32, encBuf.length)
	//   const d2 = Hash.sha256hmac(kM, c)
	//   if (toHex(d) !== toHex(d2)) {
	//     throw new Error('Invalid checksum')
	//   }
	//   const messageBuf = AESCBC.decrypt(c, kE)
	//   return [...messageBuf]

	data, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return nil, err
	}

	fromPublicKey, err := ec.ParsePubKey(data[:33])
	if err != nil {
		return nil, err
	}

	P := fromPublicKey.Mul(toPrivatKey.D)
	if P.X.Cmp(big.NewInt(0)) == 0 && P.Y.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("p equals 0")
	}

	Sbuf := P.X.Bytes()
	kEkM := c.Sha512(Sbuf)
	kE := kEkM[:32]
	kM := kEkM[32:]

	cipherText := data[33 : len(data)-32]
	mac := data[len(data)-32:]
	expectedMAC := c.Sha256HMAC(cipherText, kM)
	if !hmac.Equal(mac, expectedMAC) {
		return nil, errors.New("invalid ciphertext: HMAC mismatch")
	}
	iv := cipherText[:16]
	return ecies.AESCBCDecrypt(cipherText[16:], kE, iv)

	// Derive shared secret
	// x, _ := privateKey.Curve.ScalarMult(fromPublicKey.X, fromPublicKey.Y, privateKey.D.Bytes())
	// sharedSecret := x.Bytes()

	// Key derivation
	// keyMaterial := c.Sha512(sharedSecret)

	// keyE := keyMaterial[:32] // AES-256 key
	// keyM := keyMaterial[32:] // HMAC key

	// // Verify HMAC
	// mac := data[len(data)-32:]
	// encryptedData := data[:len(data)-32]
	// expectedMAC := c.Sha256HMAC(encryptedData, keyM)

	// if !hmac.Equal(mac, expectedMAC) {
	// 	return nil, errors.New("invalid ciphertext: HMAC mismatch")
	// }

	// // Decrypt
	// iv := make([]byte, 16) // In Bitcore, IV is usually all zeros
	// cipherText := encryptedData[33:]
	// plainText, err := ecies.AESCBCDecrypt(cipherText, keyE, iv)
	// if err != nil {
	// 	return nil, err
	// }

	// return plainText, nil
}
