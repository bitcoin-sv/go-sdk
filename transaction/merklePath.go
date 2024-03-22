package transaction

import (
	"bytes"
	"encoding/hex"
	"io"
	"sort"

	"github.com/pkg/errors"
)

type PathElement struct {
	Offset    uint64
	Hash      []byte
	Txid      bool
	Duplicate bool
}

type MerklePath struct {
	BlockHeight uint64
	Path        [][]PathElement
}

// NewMerklePath creates a new MerklePath with the given block height and path
func NewMerklePath(blockHeight uint64, path [][]PathElement) *MerklePath {
	return &MerklePath{
		BlockHeight: blockHeight,
		Path:        path,
	}
}

// NewMerklePathFromHex creates a new MerklePath with the given hex data
func NewMerklePathFromHex(hexData string) (*MerklePath, error) {
	bin, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return NewMerklePathFromBinary(bin)
}

// NewMerklePathFromBinary creates a new MerklePath with the given binary
func NewMerklePathFromBinary(bin []byte) (*MerklePath, error) {
	mp := &MerklePath{}
	_, err := mp.ReadFrom(bytes.NewReader(bin))
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// NewMerklePathFromBinary creates a new MerklePath with the given binary data
func (mp *MerklePath) ReadFrom(r io.Reader) (int64, error) {
	var bytesRead int64

	var blockHeight VarInt
	n64, err := blockHeight.ReadFrom(r)
	bytesRead += n64
	if err != nil {
		return bytesRead, err
	}

	treeHeight := make([]byte, 1)
	n, err := io.ReadFull(r, treeHeight)
	bytesRead += int64(n)
	if err != nil {
		return bytesRead, err
	}
	th := uint8(treeHeight[0])

	path := make([][]PathElement, th)

	for level := 0; level < int(th); level++ {
		var nLeavesAtThisHeight VarInt
		n64, err := nLeavesAtThisHeight.ReadFrom(r)
		bytesRead += n64
		if err != nil {
			return bytesRead, err
		}

		leaves := make([]PathElement, 0)

		for nLeavesAtThisHeight > 0 {
			var offset VarInt
			n64, err := offset.ReadFrom(r)
			bytesRead += n64
			if err != nil {
				return bytesRead, err
			}

			leaf := PathElement{
				Offset: uint64(offset),
			}
			flags := make([]byte, 1)
			n, err := io.ReadFull(r, flags)
			bytesRead += int64(n)
			if err != nil {
				return bytesRead, err
			}
			f := uint8(flags[0])

			if f&1 > 0 {
				leaf.Duplicate = true
			} else {
				if f&2 > 0 {
					leaf.Txid = true
				}
				hash := make([]byte, 32)
				n, err := io.ReadFull(r, hash)
				bytesRead += int64(n)
				if err != nil {
					return bytesRead, err
				}
				leaf.Hash = hash
			}
			leaves = append(leaves, leaf)
			nLeavesAtThisHeight--
		}
		sort.Slice(leaves, func(i, j int) bool {
			return leaves[i].Offset < leaves[j].Offset
		})
		path[level] = leaves
	}

	mp = &MerklePath{
		BlockHeight: uint64(blockHeight),
		Path:        path,
	}

	return 0, nil
}

// ToHex converts the MerklePath to a hexadecimal string representation
func (mp *MerklePath) ToHex() string {
	// This function should implement the conversion of the MerklePath to its binary representation first,
	// then convert that binary representation to a hexadecimal string.
	// Placeholder implementation below - you'll need to replace this with actual logic.
	return "hexadecimal_representation_here"
}

// ComputeRoot computes the Merkle root from a given transaction ID
func (mp *MerklePath) ComputeRoot(txid string) (string, error) {
	// Placeholder implementation. You need to implement the actual Merkle root computation based on the path and the given txid.
	// The actual computation would be significantly complex and is dependent on your specific Merkle tree structure and hashing function.
	return "", errors.New("computeRoot not implemented")
}

// // Verify checks if a given transaction ID is part of the Merkle tree at the specified block height using a chain tracker
// func (mp *MerklePath) Verify(txid string, chainTracker *ChainTracker) (bool, error) {
// 	// Placeholder for chain tracker interaction. You need to implement the verification logic here, possibly interacting with a chain tracker.
// 	// This involves computing the Merkle root and verifying it against the chain tracker's data.
// 	return false, errors.New("verify not implemented")
// }

// // Combine combines this MerklePath with another to create a compound proof
// func (mp *MerklePath) Combine(other *MerklePath) error {
// 	// Placeholder implementation. Combining two Merkle paths involves ensuring they can be combined
// 	// and then performing the combination logic.
// 	return errors.New("combine not implemented")
// }
