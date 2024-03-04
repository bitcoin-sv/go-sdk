package broadcaster

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bitcoin-sv/go-sdk/transaction"
)

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

type Arc struct {
	ApiUrl               string
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

type ArcResponse struct {
	BlockHash   string     `json:"blockHash,omitempty"`
	BlockHeight uint32     `json:"blockHeight,omitempty"`
	ExtraInfo   string     `json:"extraInfo,omitempty"`
	Status      int        `json:"status,omitempty"`
	Timestamp   time.Time  `json:"timestamp,omitempty"`
	Title       string     `json:"title,omitempty"`
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

func (a *Arc) Broadcast(t *transaction.Transaction) (interface{}, error) {
	var buf *bytes.Buffer
	for _, input := range t.Inputs {
		if input.PreviousTxScript == nil {
			buf = bytes.NewBuffer(t.Bytes())
			break
		}
	}
	if buf == nil {
		buf = bytes.NewBuffer(t.ExtendedBytes())
	}

	req, err := http.NewRequest(
		"POST",
		a.ApiUrl+"/tx",
		buf,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	if a.ApiKey != nil {
		req.Header.Set("Authorization", "Bearer "+*a.ApiKey)
	}
	if a.CallbackUrl != nil {
		req.Header.Set("X-CallbackUrl", *a.CallbackUrl)
	}
	if a.CallbackToken != nil {
		req.Header.Set("X-CallbackToken", *a.CallbackToken)
	}
	if a.FullStatusUpdates {
		req.Header.Set("X-FullStatusUpdates", "true")
	}
	if a.MaxTimeout != nil {
		req.Header.Set("X-MaxTimeout", fmt.Sprintf("%d", *a.MaxTimeout))
	}
	if a.SkipFeeValidation {
		req.Header.Set("X-SkipFeeValidation", "true")
	}
	if a.SkipScriptValidation {
		req.Header.Set("X-SkipScriptValidation", "true")
	}
	if a.SkipTxValidation {
		req.Header.Set("X-SkipTxValidation", "true")
	}
	if a.WaitForStatus != "" {
		req.Header.Set("X-WaitForStatus", string(a.WaitForStatus))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &ArcResponse{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		return nil, err
	}

	if response.Status == 200 {
		return &transaction.BroadcastSuccess{
			Status:  transaction.Success,
			Txid:    response.Txid,
			Message: response.Title,
		}, nil
	}

	return &transaction.BroadcastFailure{
		Status:      transaction.Error,
		Code:        fmt.Sprintf("%d", response.Status),
		Description: response.Title,
	}, nil
}
