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

func NewKeySharesFromBackupFormat(shares []string) (keyShares *KeyShares, error error) {
	var threshold int = 0
	var integrity string = ""
	var points []*PointInFiniteField
	for idx, share := range shares {
		shareParts := strings.Split(share, ".")
		if len(shareParts) != 4 {
			return nil, fmt.Errorf("invalid share format in share %d. Expected format: \"x.y.t.i\" - received %s", idx, share)
		}
		// convert parts to bigints
		var x *big.Int = big.NewInt(0)
		var y *big.Int = big.NewInt(0)
		// base64 decode x and y
		decodedX, err := base64.StdEncoding.DecodeString(shareParts[0])
		if err != nil {
			return nil, err
		}

		decodedY, err := base64.StdEncoding.DecodeString(shareParts[1])
		if err != nil {
			return nil, err
		}

		x.SetBytes(decodedX)
		y.SetBytes(decodedY)

		t := shareParts[2]
		i := shareParts[3]
		if t == "" {
			return nil, fmt.Errorf("threshold not found in share %d", idx)
		}
		if i == "" {
			return nil, fmt.Errorf("integrity not found in share %d", idx)
		}
		tInt, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
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
	var backupShares []string
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
