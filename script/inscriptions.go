package script

// InscriptionArgs contains the Ordinal inscription data.
type InscriptionArgs struct {
	LockingScript *Script
	Data          []byte
	ContentType   string
	EnrichedArgs  *EnrichedInscriptionArgs
}

// EnrichedInscriptionArgs contains data needed for enriched inscription
// functionality found here: https://docs.1satordinals.com/op_return.
type EnrichedInscriptionArgs struct {
	OpReturnData [][]byte
}
