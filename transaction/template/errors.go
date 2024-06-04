package template

import "errors"

var (
	ErrBadPublicKeyHash  = errors.New("invalid public key hash")
	ErrNoPrivateKey      = errors.New("private key not supplied")
	ErrBadScript         = errors.New("invalid script")
	ErrTooManySignatures = errors.New("too many signatures")
)
