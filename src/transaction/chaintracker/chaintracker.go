package chaintracker

type ChainTracker interface {
	IsValidRootForHeight(root []byte, height uint32) bool
}
