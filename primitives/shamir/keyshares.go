package primitives

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type KeyShares struct {
	Points    []*PointInFiniteField
	Threshold int
	Integrity string
}

// decodeShare decodes the share from the backup format
func decodeShare(share string) (*big.Int, *big.Int, int, string, error) {
	components := strings.Split(share, ".")
	if len(components) != 4 {
		err := fmt.Errorf("invalid share format. Expected format: \"x.y.t.i\" - received %s", share)
		return nil, nil, 0, "", err
	}
	// base64 decode x and y
	decodedX, err := base64.StdEncoding.DecodeString(components[0])
	if err != nil {
		return nil, nil, 0, "", err
	}
	decodedY, err := base64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return nil, nil, 0, "", err
	}
	var x *big.Int = big.NewInt(0).SetBytes(decodedX)
	var y *big.Int = big.NewInt(0).SetBytes(decodedY)

	t := components[2]
	if t == "" {
		return nil, nil, 0, "", fmt.Errorf("threshold not found")
	}
	i := components[3]
	if i == "" {
		return nil, nil, 0, "", fmt.Errorf("integrity not found")
	}
	tInt, err := strconv.Atoi(t)
	if err != nil {
		return nil, nil, 0, "", err
	}
	return x, y, tInt, i, nil
}

// NewKeySharesFromBackupFormat creates a new KeyShares object from a backup
func NewKeySharesFromBackupFormat(shares []string) (keyShares *KeyShares, error error) {
	var threshold int = 0
	var integrity string = ""
	points := make([]*PointInFiniteField, 0)
	for idx, share := range shares {

		x, y, tInt, i, err := decodeShare(share)
		if err != nil {
			return nil, fmt.Errorf("failed to decode share %d: %w", idx, err)
		}

		if idx != 0 && threshold != tInt {
			return nil, fmt.Errorf("threshold mismatch in share %d", idx)
		}
		if idx != 0 && integrity != i {
			return nil, fmt.Errorf("integrity mismatch in share %d", idx)
		}
		threshold = tInt
		integrity = i
		points = append(points, NewPointInFiniteField(x, y))
	}
	return NewKeyShares(points, threshold, integrity), nil
}

/**
 * @method toBackupShares
 *
 * Creates a backup of the private key by splitting it into shares.
 *
 *
 * @param threshold The number of shares which will be required to reconstruct the private key.
 * @param totalShares The number of shares to generate for distribution.
 * @returns
 */
func (k *KeyShares) ToBackupFormat() ([]string, error) {
	backupShares := make([]string, 0)
	for _, share := range k.Points {
		backupShares = append(backupShares, share.String()+"."+strconv.Itoa(k.Threshold)+"."+k.Integrity)
	}
	return backupShares, nil
}

func NewKeyShares(points []*PointInFiniteField, threshold int, integrity string) *KeyShares {
	return &KeyShares{
		Points:    points,
		Threshold: threshold,
		Integrity: integrity,
	}
}