package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	transaction "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

type SecurityLevel int

var (
	SecurityLevelSilent                  SecurityLevel = 0
	SecurityLevelEveryApp                SecurityLevel = 1
	SecurityLevelEveryAppAndCounterparty SecurityLevel = 2
)

type WalletProtocol struct {
	SecurityLevel SecurityLevel
	Protocol      string
}

type CounterpartyType int

const (
	CounterpartyTypeAnyone CounterpartyType = 0
	CounterpartyTypeSelf   CounterpartyType = 1
	CounterpartyTypeOther  CounterpartyType = 2
)

type WalletCounterparty struct {
	Type         CounterpartyType
	Counterparty *ec.PublicKey
}

type Wallet struct {
    privateKey *ec.PrivateKey
    publicKey  *ec.PublicKey
}

func NewWallet(privateKey *ec.PrivateKey) *Wallet {
    return &Wallet{
        privateKey: privateKey,
        publicKey:  privateKey.PubKey(),
    }
}

type WalletEncryptionArgs struct {
	ProtocolID       WalletProtocol
	KeyID            string
	Counterparty     WalletCounterparty
	Privileged       bool
	PrivilegedReason string
	SeekPermission   bool
	Plaintext        []byte
	Ciphertext       []byte
}

type WalletEncryptResult struct {
	Ciphertext []byte
}

type WalletDecryptResult struct {
	Plaintext []byte
}

func (w *Wallet) Encrypt(args *WalletEncryptionArgs) (*WalletEncryptResult, error) {
	if args.Counterparty.Type == CounterpartyTypeOther && args.Counterparty.Counterparty == nil {
		return nil, errors.New("counterparty public key required for other")
	}

	key := w.deriveSymmetricKey(args.ProtocolID, args.KeyID, args.Counterparty)
	ciphertext, err := encryptData(key, args.Plaintext)
	if err != nil {
		return nil, err
	}
	return &WalletEncryptResult{Ciphertext: ciphertext}, nil
}

func (w *Wallet) Decrypt(args *WalletEncryptionArgs) (*WalletDecryptResult, error) {
	if args.Counterparty.Type == CounterpartyTypeOther && args.Counterparty.Counterparty == nil {
		return nil, errors.New("counterparty public key required for other")
	}

	key := w.deriveSymmetricKey(args.ProtocolID, args.KeyID, args.Counterparty)
	plaintext, err := decryptData(key, args.Ciphertext)
	if err != nil {
		return nil, err
	}
	return &WalletDecryptResult{Plaintext: plaintext}, nil
}

func (w *Wallet) deriveSymmetricKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty) []byte {
    // If counterparty is 'anyone', use a fixed public key
    if counterparty.Type == CounterpartyTypeAnyone {
        _, fixedKey := ec.PrivateKeyFromBytes([]byte{1})
        counterparty = WalletCounterparty{
            Type:         CounterpartyTypeOther,
            Counterparty: fixedKey,
        }
    }

    // Derive both public and private keys
    derivedPublicKey := w.derivePublicKey(protocol, keyID, counterparty, false)
    derivedPrivateKey := w.derivePrivateKey(protocol, keyID, counterparty)

    // Create shared secret
    sharedSecret, _ := derivedPrivateKey.DeriveSharedSecret(derivedPublicKey)
    if sharedSecret == nil {
        return nil
    }

    // Return the x coordinate of the shared secret point
    return sharedSecret.X.Bytes()
}

func (w *Wallet) derivePublicKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty, forSelf bool) *ec.PublicKey {
    counterpartyKey := w.normalizeCounterparty(counterparty)
	var pubKey *ec.PublicKey
    if forSelf {
        privKey, _ := w.privateKey.DeriveChild(counterpartyKey, w.computeInvoiceNumber(protocol, keyID))
		if privKey != nil {
			pubKey = privKey.PubKey()
		}
    } else {
		pubKey, _ = counterpartyKey.DeriveChild(w.privateKey, w.computeInvoiceNumber(protocol, keyID))
	}
	return pubKey
}

func (w *Wallet) derivePrivateKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty) *ec.PrivateKey {
    counterpartyKey := w.normalizeCounterparty(counterparty)
    k, _ := w.privateKey.DeriveChild(counterpartyKey, w.computeInvoiceNumber(protocol, keyID))
	return k
}

func (w *Wallet) normalizeCounterparty(counterparty WalletCounterparty) *ec.PublicKey {
    switch counterparty.Type {
    case CounterpartyTypeSelf:
        return w.privateKey.PubKey()
    case CounterpartyTypeOther:
        return counterparty.Counterparty
    case CounterpartyTypeAnyone:
        _, pub := ec.PrivateKeyFromBytes([]byte{1})
		return pub
    default:
        return nil
    }
}

func (w *Wallet) computeInvoiceNumber(protocol WalletProtocol, keyID string) string {
    // Validate protocol security level
    if protocol.SecurityLevel < 0 || protocol.SecurityLevel > 2 {
        return ""
    }

    // Validate protocol name
    protocolName := strings.ToLower(strings.TrimSpace(protocol.Protocol))
    if len(protocolName) > 400 || len(protocolName) < 5 {
        return ""
    }
    if strings.Contains(protocolName, "  ") {
        return ""
    }
    if !regexp.MustCompile(`^[a-z0-9 ]+$`).MatchString(protocolName) {
        return ""
    }
    if strings.HasSuffix(protocolName, " protocol") {
        return ""
    }

    // Validate key ID
    if len(keyID) > 800 || len(keyID) < 1 {
        return ""
    }

    return fmt.Sprintf("%d-%s-%s", protocol.SecurityLevel, protocolName, keyID)
}

func encryptData(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

type GetPublicKeyArgs struct {
	WalletEncryptionArgs
	IdentityKey bool
	ForSelf     bool
}

type GetPublicKeyResult struct {
	PublicKey *ec.PublicKey `json:"publicKey"`
}

func (w *Wallet) GetPublicKey(args *GetPublicKeyArgs, originator string) GetPublicKeyResult {
	return GetPublicKeyResult{
		PublicKey: &ec.PublicKey{},
	}
}

type CreateSignatureArgs struct {
	WalletEncryptionArgs
	Data               []byte
	DashToDirectlySign []byte
}

type CreateSignatureResult struct {
	Signature ec.Signature
}

type SignOutputs transaction.Flag

var (
	SignOutputsAll    SignOutputs = SignOutputs(sighash.All)
	SignOutputsNone   SignOutputs = SignOutputs(sighash.None)
	SignOutputsSingle SignOutputs = SignOutputs(sighash.Single)
)

func (w *Wallet) CreateSignature(args *CreateSignatureArgs, originator string) CreateSignatureResult {
	// a := sighash.All
	return CreateSignatureResult{
		Signature: ec.Signature{},
	}
}
