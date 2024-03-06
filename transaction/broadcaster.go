package transaction

type BroadcastStatus string

var (
	Success BroadcastStatus = "success"
	Error   BroadcastStatus = "error"
)

type BroadcastSuccess struct {
	Status  BroadcastStatus `json:"status"`
	Txid    string          `json:"txid"`
	Message string          `json:"message"`
}

type BroadcastFailure struct {
	Status      BroadcastStatus `json:"status"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
}

type Broadcaster interface {
	Broadcast(tx *Tx) (interface{}, error)
}

func (t *Tx) Broadcast(b Broadcaster) (interface{}, error) {
	return b.Broadcast(t)
}
