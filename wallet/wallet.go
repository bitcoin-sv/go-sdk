package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	sighash "github.com/bsv-blockchain/go-sdk/transaction/sighash"
	transaction "github.com/bsv-blockchain/go-sdk/transaction/sighash"
	"io"
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
	keyDeriver *KeyDeriver
}

func NewWallet(privateKey *ec.PrivateKey) *Wallet {
	return &Wallet{
		privateKey: privateKey,
		publicKey:  privateKey.PubKey(),
		keyDeriver: NewKeyDeriver(privateKey),
	}
}

type WalletEncryptionArgs struct {
	ProtocolID       WalletProtocol
	KeyID            string
	Counterparty     WalletCounterparty
	Privileged       bool
	PrivilegedReason string
	SeekPermission   bool
}

type WalletEncryptArgs struct {
	WalletEncryptionArgs
	Plaintext []byte
}

type WalletDecryptArgs struct {
	WalletEncryptionArgs
	Ciphertext []byte
}

type WalletEncryptResult struct {
	Ciphertext []byte
}

type WalletDecryptResult struct {
	Plaintext []byte
}

func (w *Wallet) Encrypt(args *WalletEncryptArgs) (*WalletEncryptResult, error) {
	if args.Counterparty.Type == CounterpartyTypeOther && args.Counterparty.Counterparty == nil {
		return nil, errors.New("counterparty public key required for other")
	}

	key, err := w.keyDeriver.DeriveSymmetricKey(args.ProtocolID, args.KeyID, args.Counterparty)
	if err != nil {
		return nil, fmt.Errorf("failed to derive symmetric key: %w", err)
	}

	ciphertext, err := encryptData(key, args.Plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	return &WalletEncryptResult{Ciphertext: ciphertext}, nil
}

func (w *Wallet) Decrypt(args *WalletDecryptArgs) (*WalletDecryptResult, error) {
	if args.Counterparty.Type == CounterpartyTypeOther && args.Counterparty.Counterparty == nil {
		return nil, errors.New("counterparty public key required for other")
	}

	key, err := w.keyDeriver.DeriveSymmetricKey(args.ProtocolID, args.KeyID, args.Counterparty)
	if err != nil {
		return nil, fmt.Errorf("failed to derive symmetric key: %w", err)
	}

	plaintext, err := decryptData(key, args.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return &WalletDecryptResult{Plaintext: plaintext}, nil
}

func encryptData(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: expected at least %d bytes, got %d", nonceSize, len(ciphertext))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
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

func (w *Wallet) CreateSignature(args *CreateSignatureArgs, originator string) (*CreateSignatureResult, error) {
	if len(args.Data) == 0 && len(args.DashToDirectlySign) == 0 {
		return nil, fmt.Errorf("args.data or args.hashToDirectlySign must be valid")
	}

	// Get hash to sign
	var hash []byte
	if len(args.DashToDirectlySign) > 0 {
		hash = args.DashToDirectlySign
	} else {
		sum := sha256.Sum256(args.Data)
		hash = sum[:]
	}

	// Derive private key
	privKey, err := w.keyDeriver.DerivePrivateKey(
		args.ProtocolID,
		args.KeyID,
		args.Counterparty,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	// Create signature
	signature, err := privKey.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	return &CreateSignatureResult{
		Signature: *signature,
	}, nil
}

type VerifySignatureArgs struct {
	WalletEncryptionArgs
	Data                []byte
	DashToDirectlyVerify []byte
	Signature           ec.Signature
	ForSelf             bool
}

type VerifySignatureResult struct {
	Valid bool
}

func (w *Wallet) VerifySignature(args *VerifySignatureArgs, originator string) (*VerifySignatureResult, error) {
	if len(args.Data) == 0 && len(args.DashToDirectlyVerify) == 0 {
		return nil, fmt.Errorf("args.data or args.hashToDirectlyVerify must be valid")
	}

	// Get hash to verify
	var hash []byte
	if len(args.DashToDirectlyVerify) > 0 {
		hash = args.DashToDirectlyVerify
	} else {
		sum := sha256.Sum256(args.Data)
		hash = sum[:]
	}

	// Derive public key
	pubKey, err := w.keyDeriver.DerivePublicKey(
		args.ProtocolID,
		args.KeyID,
		args.Counterparty,
		args.ForSelf,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %w", err)
	}

	// Verify signature
	valid := args.Signature.Verify(hash, pubKey)
	if !valid {
		return nil, fmt.Errorf("signature is not valid")
	}

	return &VerifySignatureResult{
		Valid: valid,
	}, nil
}
