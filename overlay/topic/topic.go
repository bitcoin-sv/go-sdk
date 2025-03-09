package topic

import "github.com/bitcoin-sv/go-sdk/overlay"

type TaggedBEEF struct {
	Beef   []byte
	Topics []string
}

type Facilitator interface {
	Send(url string, taggedBEEF *TaggedBEEF) (*overlay.Steak, error)
}
