package transaction

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/sighash"
)

// UnlockerParams params used for unlocking an input with a `bt.Unlocker`.
type UnlockerParams struct {
	// InputIdx the input to be unlocked. [DEFAULT 0]
	InputIdx uint32
	// SigHashFlags the be applied [DEFAULT ALL|FORKID]
	SigHashFlags sighash.Flag
	// TODO: add previous tx script and sats here instead of in
	// input (and potentially remove from input) - see issue #143
}

// Unlocker interface to allow custom implementations of different unlocking mechanisms.
// Implement the Unlocker function as shown in LocalUnlocker, for example.
type Unlocker interface {
	UnlockingScript(ctx context.Context, tx *Transaction, up UnlockerParams) (uscript *script.Script, err error)
}

// UnlockerGetter interfaces getting an unlocker for a given output/locking script.
type UnlockerGetter interface {
	Unlocker(ctx context.Context, lockingScript *script.Script) (Unlocker, error)
}