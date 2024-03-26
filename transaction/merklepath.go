package transaction

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/bitcoin-sv/go-sdk/crypto"
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
	if len(txid) == 0 {
		for _, l := range mp.Path[0] {
			if len(l.Hash) > 0 {
				txid = hex.EncodeToString(l.Hash)
				break
			}
		}
	}
	if len(mp.Path) == 1 {
		// if there is only one txid in the block then the root is the txid.
		if len(mp.Path[0]) == 1 {
			return txid, nil
		}
	}
	// Find the index of the txid at the lowest level of the Merkle tree
	var index uint64
	txidFound := false
	for _, l := range mp.Path[0] {
		if hex.EncodeToString(l.Hash) == txid {
			txidFound = true
			index = l.Offset
			break
		}
	}
	if !txidFound {
		return "", errors.New("the BUMP does not contain the txid: " + txid)
	}

	// Calculate the root using the index as a way to determine which direction to concatenate.
	workingHash, err := hex.DecodeString(txid)
	if err != nil {
		return "", err
	}
	workingHash = ReverseBytes(workingHash)
	for height, leaves := range mp.Path {
		offset := (index >> height) ^ 1
		var leafAtThisLevel PathElement
		offsetFound := false
		for _, l := range leaves {
			if l.Offset == offset {
				offsetFound = true
				leafAtThisLevel = l
				break
			}
		}
		if !offsetFound {
			return "", fmt.Errorf("we do not have a hash for this index at height: %v", height)
		}

		var digest []byte
		if leafAtThisLevel.Duplicate {
			digest = append(workingHash, workingHash...)
		} else {
			leafBytes := ReverseBytes(leafAtThisLevel.Hash)
			if (offset % 2) != 0 {
				digest = append(workingHash, leafBytes...)
			} else {
				digest = append(leafBytes, workingHash...)
			}
		}
		workingHash = crypto.Sha256d(digest)
	}
	return hex.EncodeToString(ReverseBytes(workingHash)), nil
}

// Verify checks if a given transaction ID is part of the Merkle tree at the specified block height using a chain tracker
func (mp *MerklePath) Verify(txid string, ct chaintracker.ChainTracker) (bool, error) {
	root, err := mp.ComputeRoot(txid)
	if err != nil {
		return false, err
	}
	rootBytes, err := hex.DecodeString(root)
	if err != nil {
		return false, err
	}
	return ct.IsValidRootForHeight(rootBytes, mp.BlockHeight), nil
}

// Combine combines this MerklePath with another to create a compound proof
// const root1 = this.computeRoot()
//
//	const root2 = other.computeRoot()
//	if (root1 !== root2) {
//	  throw Error('You cannot combine paths which do not have the same root.')
//	}
//	const combinedPath = []
//	for (let h = 0; h < this.path.length; h++) {
//	  combinedPath.push([])
//	  for (let l = 0; l < this.path[h].length; l++) {
//	    combinedPath[h].push(this.path[h][l])
//	  }
//	  for (let l = 0; l < other.path[h].length; l++) {
//	    if (!(combinedPath[h].find(leaf => leaf.offset === other.path[h][l].offset) as boolean)) {
//	      combinedPath[h].push(other.path[h][l])
//	    } else {
//	      // Ensure that any elements which appear in both are not downgraded to a non txid.
//	      if (other.path[h][l]?.txid) {
//	        combinedPath[h].find(leaf => leaf.offset === other.path[h][l]).txid = true
//	      }
//	    }
//	  }
//	}
//	this.path = combinedPath
func (mp *MerklePath) Combine(other *MerklePath) error {
	if mp.BlockHeight != other.BlockHeight {
		return errors.New("cannot combine MerklePaths with different block heights")
	}

	root1, err := mp.ComputeRoot("")
}
