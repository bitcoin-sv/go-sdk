package chaintracker

type ChainTracker interface {
	IsValidRootForHeight(root []byte, height uint64) bool
}
