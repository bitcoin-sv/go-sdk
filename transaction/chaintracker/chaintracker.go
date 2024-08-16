package chaintracker

import "github.com/bitcoin-sv/go-sdk/chainhash"

type ChainTracker interface {
	IsValidRootForHeight(root *chainhash.Hash, height uint32) bool
}
