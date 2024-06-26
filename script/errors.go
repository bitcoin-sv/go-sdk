package script

import "github.com/pkg/errors"

// Sentinel errors raised by data ops.
var (
	ErrDataTooBig            = errors.New("data too big")
	ErrDataTooSmall          = errors.New("not enough data")
	ErrPartTooBig            = errors.New("part too big")
	ErrScriptIndexOutOfRange = errors.New("script index out of range")
)

// Sentinel errors raised by addresses.
var (
	ErrInvalidAddressLength = errors.New("invalid address length")
	ErrUnsupportedAddress   = errors.New("address not supported")
)

// Sentinel errors raised by inscriptions.
var (
	ErrP2PKHInscriptionNotFound = errors.New("no P2PKH inscription found")
)

// Sentinel errors raised through encoding.
var (
	ErrEncodingBadChar         = errors.New("bad char")
	ErrEncodingTooLong         = errors.New("too long")
	ErrEncodingInvalidVersion  = errors.New("not version 0 of 6f")
	ErrEncodingInvalidChecksum = errors.New("invalid checksum")
	ErrEncodingChecksumFailed  = errors.New("checksum failed")
	ErrTextNoBIP76             = errors.New("text did not match the bip276 format")
)

// Sentinel errors raised by the package.
var (
	ErrInvalidPKLen      = errors.New("invalid public key length")
	ErrInvalidOpCode     = errors.New("invalid opcode data")
	ErrEmptyScript       = errors.New("script is empty")
	ErrNotP2PKH          = errors.New("not a P2PKH")
	ErrInvalidOpcodeType = errors.New("use AppendPushData for push data funcs")
)
