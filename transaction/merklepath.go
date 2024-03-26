package transaction

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/bitcoin-sv/go-sdk/transaction/chaintracker"
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
func NewMerklePathFromBinary(bytes []byte) (*MerklePath, error) {
	if len(bytes) < 37 {
		return nil, errors.New("BUMP bytes do not contain enough data to be valid")
	}
	bump := &MerklePath{}

	// first bytes are the block height.
	var skip int
	index, size := NewVarIntFromBytes(bytes[skip:])
	skip += size
	bump.BlockHeight = uint64(index)

	// Next byte is the tree height.
	treeHeight := uint(bytes[skip])
	skip++

	// We expect tree height levels.
	bump.Path = make([][]PathElement, treeHeight)

	for lv := uint(0); lv < treeHeight; lv++ {
		// For each level we parse a bunch of nLeaves.
		n, size := NewVarIntFromBytes(bytes[skip:])
		skip += size
		nLeavesAtThisHeight := uint64(n)
		if nLeavesAtThisHeight == 0 {
			return nil, errors.New("There are no leaves at height: " + fmt.Sprint(lv) + " which makes this invalid")
		}
		bump.Path[lv] = make([]PathElement, nLeavesAtThisHeight)
		for lf := uint64(0); lf < nLeavesAtThisHeight; lf++ {
			// For each leaf we parse the offset, hash, txid and duplicate.
			offset, size := NewVarIntFromBytes(bytes[skip:])
			skip += size
			var l PathElement
			o := uint64(offset)
			l.Offset = o
			flags := bytes[skip]
			skip++
			dup := flags&1 > 0
			txid := flags&2 > 0
			if dup {
				l.Duplicate = dup
			} else {
				if len(bytes) < skip+32 {
					return nil, errors.New("BUMP bytes do not contain enough data to be valid")
				}
				h := bytes[skip : skip+32]
				l.Hash = h
				skip += 32
			}
			if txid {
				l.Txid = txid
			}
			bump.Path[lv][lf] = l
		}
	}

	// Sort each of the levels by the offset for consistency.
	for _, level := range bump.Path {
		sort.Slice(level, func(i, j int) bool {
			return level[i].Offset < level[j].Offset
		})
	}

	return bump, nil
}

// Bytes encodes a BUMP as a slice of bytes. BUMP Binary Format according to BRC-74 https://brc.dev/74
func (mp *MerklePath) Bytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, VarInt(mp.BlockHeight).Bytes()...)
	treeHeight := len(mp.Path)
	bytes = append(bytes, byte(treeHeight))
	for level := 0; level < treeHeight; level++ {
		nLeaves := len(mp.Path[level])
		bytes = append(bytes, VarInt(nLeaves).Bytes()...)
		for _, leaf := range mp.Path[level] {
			bytes = append(bytes, VarInt(leaf.Offset).Bytes()...)
			flags := byte(0)
			if leaf.Duplicate {
				flags |= 1
			}
			if leaf.Txid {
				flags |= 2
			}
			bytes = append(bytes, flags)
			if (flags & 1) == 0 {
				bytes = append(bytes, ReverseBytes(leaf.Hash)...)
			}
		}
	}
	return bytes
}

// ToHex converts the MerklePath to a hexadecimal string representation
func (mp *MerklePath) ToHex() string {
	return hex.EncodeToString(mp.Bytes())
}

// ComputeRoot computes the Merkle root from a given transaction ID
func (mp *MerklePath) ComputeRoot(txid string) (string, error) {
	// Placeholder implementation. You need to implement the actual Merkle root computation based on the path and the given txid.
	// The actual computation would be significantly complex and is dependent on your specific Merkle tree structure and hashing function.
	return "", errors.New("computeRoot not implemented")
}

// Verify checks if a given transaction ID is part of the Merkle tree at the specified block height using a chain tracker
func (mp *MerklePath) Verify(txid string, chainTracker chaintracker.ChainTracker) (bool, error) {
	// Placeholder for chain tracker interaction. You need to implement the verification logic here, possibly interacting with a chain tracker.
	// This involves computing the Merkle root and verifying it against the chain tracker's data.
	return false, errors.New("verify not implemented")
}

// Combine combines this MerklePath with another to create a compound proof
func (mp *MerklePath) Combine(other *MerklePath) error {
	// Placeholder implementation. Combining two Merkle paths involves ensuring they can be combined
	// and then performing the combination logic.
	return errors.New("combine not implemented")
}
