package transaction

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ArcResponse struct {
	BlockHash   string     `json:"blockHash,omitempty"`
	BlockHeight uint32     `json:"blockHeight,omitempty"`
	ExtraInfo   string     `json:"extraInfo,omitempty"`
	Status      int        `json:"status,omitempty"`
	Timestamp   time.Time  `json:"timestamp,omitempty"`
	TxStatus    *ArcStatus `json:"txStatus,omitempty"`
	Instance    *string    `json:"instance,omitempty"`
	Txid        string     `json:"txid,omitempty"`
	Detail      *string    `json:"detail,omitempty"`
}

func (ts ArcResponse) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

func (f *ArcResponse) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &f)
}

type ArcStatus string

const (
	RECEIVED             ArcStatus = "2"
	STORED               ArcStatus = "3"
	ANNOUNCED_TO_NETWORK ArcStatus = "4"
	REQUESTED_BY_NETWORK ArcStatus = "5"
	SENT_TO_NETWORK      ArcStatus = "6"
	ACCEPTED_BY_NETWORK  ArcStatus = "7"
	SEEN_ON_NETWORK      ArcStatus = "8"
)

type ArcOptions struct {
	ApiKey               *string
	CallbackUrl          *string
	CallbackToken        *string
	FullStatusUpdates    bool
	MaxTimeout           *int
	SkipFeeValidation    bool
	SkipScriptValidation bool
	SkipTxValidation     bool
	WaitForStatus        ArcStatus
}

func (t *Transaction) BroadcastToArc(arcUrl string, options *ArcOptions) (response *ArcResponse, err error) {
	buf := bytes.NewBuffer(t.ToEF())
	req, err := http.NewRequest(
		"POST",
		arcUrl,
		buf,
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	if options.ApiKey != nil {
		req.Header.Set("Authorization", "Bearer "+*options.ApiKey)
	}
	if options.CallbackUrl != nil {
		req.Header.Set("X-CallbackUrl", *options.CallbackUrl)
	}
	if options.CallbackToken != nil {
		req.Header.Set("X-CallbackToken", *options.CallbackToken)
	}
	if options.FullStatusUpdates {
		req.Header.Set("X-FullStatusUpdates", "true")
	}
	if options.MaxTimeout != nil {
		req.Header.Set("X-MaxTimeout", fmt.Sprintf("%d", *options.MaxTimeout))
	}
	if options.SkipFeeValidation {
		req.Header.Set("X-SkipFeeValidation", "true")
	}
	if options.SkipScriptValidation {
		req.Header.Set("X-SkipScriptValidation", "true")
	}
	if options.SkipTxValidation {
		req.Header.Set("X-SkipTxValidation", "true")
	}
	if options.WaitForStatus != "" {
		req.Header.Set("X-WaitForStatus", string(options.WaitForStatus))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response = &ArcResponse{}
	err = json.Unmarshal(msg, &response)
	return
}
