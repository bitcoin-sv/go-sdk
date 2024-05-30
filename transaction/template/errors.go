package template

import "errors"

var (
	ErrInvalidPublicKeyHash = errors.New("invalid public key hash")
	ErrMisingPrivateKey     = errors.New("private key not supplied")
)
