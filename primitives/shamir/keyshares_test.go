package primitives

import (
	"encoding/base64"
	"math/big"
	"testing"
)

func TestNewKeyshares(t *testing.T) {
	points := []*PointInFiniteField{
		NewPointInFiniteField(big.NewInt(0), big.NewInt(1)),
		NewPointInFiniteField(big.NewInt(1), big.NewInt(2)),
		NewPointInFiniteField(big.NewInt(2), big.NewInt(3)),
	}
	threshold := 3
	integrity := base64.StdEncoding.EncodeToString([]byte("integrity"))
	keyShares := NewKeyShares(points, threshold, integrity)
	if keyShares == nil {
		t.Errorf("Failed to create new key shares")
	}

	// test backup format
	backup, err := keyShares.ToBackupFormat()
	if err != nil {
		t.Errorf("Failed to convert key shares to backup format: %v", err)
	}
	if len(backup) != 3 {
		t.Errorf("Expected 3 shares, got %d", len(backup))
	}

	newShares, err := NewKeySharesFromBackupFormat(backup)
	if err != nil {
		t.Errorf("Failed to create key shares from backup format: %v", err)
	}

	if keyShares == nil {
		t.Errorf("Failed to create new key shares")
		return
	}

	if newShares.Threshold != keyShares.Threshold {
		t.Errorf("Threshold mismatch. Expected %d, got %d", keyShares.Threshold, newShares.Threshold)
	}
}
