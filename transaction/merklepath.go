package transaction

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/bitcoin-sv/go-sdk/crypto"
	"github.com/bitcoin-sv/go-sdk/transaction/chaintracker"
	"github.com/bitcoin-sv/go-sdk/util"
	"github.com/pkg/errors"
)

type PathElement struct {
	Offset    uint64
	Hash      util.ByteStringLE
	Txid      bool
	Duplicate bool
}

type MerklePath struct {
	BlockHeight uint32
	Path        [][]*PathElement
}

// NewMerklePath creates a new MerklePath with the given block height and path
func NewMerklePath(blockHeight uint32, path [][]*PathElement) *MerklePath {
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

	var skip int
	index, size := NewVarIntFromBytes(bytes[skip:])
	skip += size
	bump.BlockHeight = uint32(index)

	// Next byte is the tree height.
	treeHeight := uint(bytes[skip])
	skip++

	// We expect tree height levels.
	bump.Path = make([][]*PathElement, treeHeight)

	for lv := uint(0); lv < treeHeight; lv++ {
		// For each level we parse a bunch of nLeaves.
		if len(bytes) <= skip {
			return nil, errors.New("Malformed BUMP")
		}
		n, size := NewVarIntFromBytes(bytes[skip:])
		skip += size
		nLeavesAtThisHeight := uint64(n)
		if nLeavesAtThisHeight == 0 {
			return nil, errors.New("There are no leaves at height: " + fmt.Sprint(lv) + " which makes this invalid")
		}
		bump.Path[lv] = make([]*PathElement, nLeavesAtThisHeight)
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
				l.Hash = bytes[skip : skip+32]
				skip += 32
			}
			if txid {
				l.Txid = txid
			}
			bump.Path[lv][lf] = &l
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
	bytes := VarInt(mp.BlockHeight).Bytes()
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
				bytes = append(bytes, leaf.Hash...)
			}
		}
	}
	return bytes
}

// ToHex converts the MerklePath to a hexadecimal string representation
func (mp *MerklePath) ToHex() string {
	return hex.EncodeToString(mp.Bytes())
}

func (mp *MerklePath) ComputeRoot(txid *string) (string, error) {
	var txidLE *[]byte
	if txid != nil {
		txidBytes, err := hex.DecodeString(*txid)
		if err != nil {
			return "", err
		}
		txidBytes = util.ReverseBytes(txidBytes)
		txidLE = &txidBytes
	}
	root, err := mp.ComputeRootBin(txidLE)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(util.ReverseBytes(root)), nil
}

// ComputeRoot computes the Merkle root from a given transaction ID
func (mp *MerklePath) ComputeRootBin(txidLE *[]byte) ([]byte, error) {
	if txidLE == nil {
		for _, l := range mp.Path[0] {
			if len(l.Hash) > 0 {
				t := []byte(l.Hash)
				txidLE = &t
				break
			}
		}
	}
	if len(mp.Path) == 1 {
		// if there is only one txid in the block then the root is the txid.
		if len(mp.Path[0]) == 1 {
			return *txidLE, nil
		}
	}

	// Find the index of the txid at the lowest level of the Merkle tree
	var txLeaf *PathElement
	for _, l := range mp.Path[0] {
		if bytes.Equal(l.Hash, *txidLE) {
			txLeaf = l
			break
		}
	}
	if txLeaf == nil {
		return nil, fmt.Errorf("the BUMP does not contain the txid: %x", *txidLE)
	}

	// Calculate the root using the index as a way to determine which direction to concatenate.
	workingHash := txLeaf.Hash
	index := txLeaf.Offset

	for height, leaves := range mp.Path {
		offset := (index >> height) ^ 1
		var leafAtThisLevel PathElement
		offsetFound := false
		for _, l := range leaves {
			if l.Offset == offset {
				offsetFound = true
				leafAtThisLevel = *l
				break
			}
		}
		if !offsetFound {
			return nil, fmt.Errorf("we do not have a hash for this index at height: %v", height)
		}

		var digest []byte
		if leafAtThisLevel.Duplicate {
			digest = append(workingHash, workingHash...)
		} else {
			leafBytes := leafAtThisLevel.Hash
			if (offset % 2) != 0 {
				digest = append(workingHash, leafBytes...)
			} else {
				digest = append(leafBytes, workingHash...)
			}
		}
		workingHash = crypto.Sha256d(digest)
	}
	return workingHash, nil
}

// Verify checks if a given transaction ID is part of the Merkle tree at the specified block height using a chain tracker
func (mp *MerklePath) Verify(txid string, ct chaintracker.ChainTracker) (bool, error) {
	root, err := mp.ComputeRoot(&txid)
	if err != nil {
		return false, err
	}
	rootBytes, err := hex.DecodeString(root)
	if err != nil {
		return false, err
	}
	rootBytes = util.ReverseBytes(rootBytes)
	return ct.IsValidRootForHeight(rootBytes, mp.BlockHeight), nil
}

func (m *MerklePath) Combine(other *MerklePath) (err error) {
	if m.BlockHeight != other.BlockHeight {
		return errors.New("cannot combine MerklePaths with different block heights")
	}

	root1, err := m.ComputeRoot(nil)
	if err != nil {
		return err
	}
	root2, err := other.ComputeRoot(nil)
	if err != nil {
		return err
	}

	if root1 != root2 {
		return errors.New("cannot combine MerklePaths with different roots")
	}

	combinedPath := make([][]*PathElement, len(m.Path))
	for h := 0; h < len(m.Path); h++ {
		for l := 0; l < len(m.Path[h]); l++ {
			combinedPath[h] = append(combinedPath[h], m.Path[h][l])
		}
		for l := 0; l < len(other.Path[h]); l++ {
			var found *PathElement
			for _, leaf := range combinedPath[h] {
				if leaf.Offset == other.Path[h][l].Offset {
					found = other.Path[h][l]
					break
				}
			}
			if found == nil {
				combinedPath[h] = append(combinedPath[h], other.Path[h][l])
			} else {
				for _, leaf := range combinedPath[h] {
					if leaf.Offset == other.Path[h][l].Offset {
						leaf.Txid = true
						break
					}
				}
			}
		}
	}
	m.Path = combinedPath
	return
}
