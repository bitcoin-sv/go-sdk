package transaction

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/script"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
)

// UnlockerParams params used for unlocking an input with a `bt.Unlocker`.
type UnlockerParams struct {
	// InputIdx the input to be unlocked. [DEFAULT 0]
	InputIdx uint32
	// SigHashFlags the be applied [DEFAULT ALL|FORKID]
	SigHashFlags sighash.Flag
}

// Unlocker interface to allow custom implementations of different unlocking mechanisms.
// Implement the Unlocker function as shown in LocalUnlocker, for example.
type Unlocker interface {
	Sign(ctx context.Context, tx *Transaction, up UnlockerParams) (uscript *script.Script, err error)
	EstimateSize() int
}

// // UnlockerGetter interfaces getting an unlocker for a given output/locking script.
// type UnlockerGetter interface {
// 	Unlocker(ctx context.Context, lockingScript *script.Script) (Unlocker, error)
// }
