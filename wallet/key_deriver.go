package wallet

import (
	"fmt"
	"regexp"
	"strings"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

type KeyDeriver struct {
	privateKey *ec.PrivateKey
	publicKey  *ec.PublicKey
}

func NewKeyDeriver(privateKey *ec.PrivateKey) *KeyDeriver {
	return &KeyDeriver{
		privateKey: privateKey,
		publicKey:  privateKey.PubKey(),
	}
}

func (kd *KeyDeriver) DeriveSymmetricKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty) []byte {
	// If counterparty is 'anyone', use a fixed public key
	if counterparty.Type == CounterpartyTypeAnyone {
		_, fixedKey := ec.PrivateKeyFromBytes([]byte{1})
		counterparty = WalletCounterparty{
			Type:         CounterpartyTypeOther,
			Counterparty: fixedKey,
		}
	}

	// Derive both public and private keys
	derivedPublicKey := kd.DerivePublicKey(protocol, keyID, counterparty, false)
	derivedPrivateKey := kd.DerivePrivateKey(protocol, keyID, counterparty)

	// Create shared secret
	sharedSecret, _ := derivedPrivateKey.DeriveSharedSecret(derivedPublicKey)
	if sharedSecret == nil {
		return nil
	}

	// Return the x coordinate of the shared secret point
	return sharedSecret.X.Bytes()
}

func (kd *KeyDeriver) DerivePublicKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty, forSelf bool) *ec.PublicKey {
	counterpartyKey := kd.NormalizeCounterparty(counterparty)
	var pubKey *ec.PublicKey
	if forSelf {
		privKey, _ := kd.privateKey.DeriveChild(counterpartyKey, kd.ComputeInvoiceNumber(protocol, keyID))
		if privKey != nil {
			pubKey = privKey.PubKey()
		}
	} else {
		pubKey, _ = counterpartyKey.DeriveChild(kd.privateKey, kd.ComputeInvoiceNumber(protocol, keyID))
	}
	return pubKey
}

func (kd *KeyDeriver) DerivePrivateKey(protocol WalletProtocol, keyID string, counterparty WalletCounterparty) *ec.PrivateKey {
	counterpartyKey := kd.NormalizeCounterparty(counterparty)
	k, _ := kd.privateKey.DeriveChild(counterpartyKey, kd.ComputeInvoiceNumber(protocol, keyID))
	return k
}

func (kd *KeyDeriver) NormalizeCounterparty(counterparty WalletCounterparty) *ec.PublicKey {
	switch counterparty.Type {
	case CounterpartyTypeSelf:
		return kd.privateKey.PubKey()
	case CounterpartyTypeOther:
		return counterparty.Counterparty
	case CounterpartyTypeAnyone:
		_, pub := ec.PrivateKeyFromBytes([]byte{1})
		return pub
	default:
		return nil
	}
}

func (kd *KeyDeriver) ComputeInvoiceNumber(protocol WalletProtocol, keyID string) string {
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
