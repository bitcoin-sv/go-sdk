package lookup

import "time"

type Facilitator interface {
	Lookup(url string, question LookupQuestion, timeout time.Duration) (LookupAnswer, error)
}
