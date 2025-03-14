package broadcaster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bsv-blockchain/go-sdk/transaction"
)

type WOCNetwork string

var (
	WOCMainnet WOCNetwork = "main"
	WOCTestnet WOCNetwork = "test"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WhatsOnChain struct {
	Network WOCNetwork
	ApiKey  string
	Client  HTTPClient
}

func (b *WhatsOnChain) Broadcast(t *transaction.Transaction) (
	*transaction.BroadcastSuccess,
	*transaction.BroadcastFailure,
) {
	if t == nil {
		return nil, &transaction.BroadcastFailure{
			Code:        "500",
			Description: "nil transaction",
		}
	}

	if b.Client == nil {
		b.Client = http.DefaultClient
	}

	bodyMap := map[string]interface{}{
		"txhex": t.Hex(),
	}
	if body, err := json.Marshal(bodyMap); err != nil {
		return nil, &transaction.BroadcastFailure{
			Code:        "500",
			Description: err.Error(),
		}
	} else {
		url := fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s/tx/raw", b.Network)
		ctx := context.Background()
		req, err := http.NewRequestWithContext(
			ctx,
			"POST",
			url,
			bytes.NewBuffer(body),
		)
		if err != nil {
			return nil, &transaction.BroadcastFailure{
				Code:        "500",
				Description: err.Error(),
			}
		}
		req.Header.Set("Content-Type", "application/json")
		if b.ApiKey != "" {
			req.Header.Set("Authorization", "Bearer "+b.ApiKey)
		}

		if resp, err := b.Client.Do(req); err != nil {
			return nil, &transaction.BroadcastFailure{
				Code:        "500",
				Description: err.Error(),
			}
		} else {
			defer resp.Body.Close() //nolint:errcheck // standard http client pattern
			if resp.StatusCode != 200 {
				if body, err := io.ReadAll(resp.Body); err != nil {
					return nil, &transaction.BroadcastFailure{
						Code:        fmt.Sprintf("%d", resp.StatusCode),
						Description: "unknown error",
					}
				} else {
					return nil, &transaction.BroadcastFailure{
						Code:        fmt.Sprintf("%d", resp.StatusCode),
						Description: string(body),
					}
				}
			} else {
				return &transaction.BroadcastSuccess{
					Txid: t.TxID().String(),
				}, nil
			}
		}
	}
}
