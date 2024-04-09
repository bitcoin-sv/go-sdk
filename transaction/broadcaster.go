package transaction

type BroadcastStatus string

var (
	Success BroadcastStatus = "success"
	Error   BroadcastStatus = "error"
)

type BroadcastSuccess struct {
	Txid    string `json:"txid"`
	Message string `json:"message"`
}

type BroadcastFailure struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *BroadcastFailure) Error() string {
	return e.Description
}

type Broadcaster interface {
	Broadcast(tx *Tx) (*BroadcastSuccess, *BroadcastFailure)
}

func (t *Tx) Broadcast(b Broadcaster) (*BroadcastSuccess, *BroadcastFailure) {
	return b.Broadcast(t)
}
