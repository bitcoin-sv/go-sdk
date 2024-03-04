package transaction

type PathElement struct {
	Offset    uint32
	Hash      []byte
	Txid      bool
	Duplicate bool
}

type MerklePath struct {
	BlockHeight uint32
	Path        [][]PathElement
}
