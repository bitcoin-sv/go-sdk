package chaintracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/go-sdk/chainhash"
)

type Network string

type BlockHeader struct {
	Hash       *chainhash.Hash `json:"hash"`
	Height     uint32          `json:"height"`
	Version    uint32          `json:"version"`
	MerkleRoot *chainhash.Hash `json:"merkleroot"`
	Time       uint32          `json:"time"`
	Nonce      uint32          `json:"nonce"`
	Bits       string          `json:"bits"`
	PrevHash   *chainhash.Hash `json:"previousblockhash"`
}

var (
	MainNet Network = "main"
	TestNet Network = "test"
)

type WhatsOnChain struct {
	Network Network
	ApiKey  string
}

func NewWhatsOnChain(network Network, apiKey string) *WhatsOnChain {
	return &WhatsOnChain{
		Network: network,
		ApiKey:  apiKey,
	}
}

func (w *WhatsOnChain) GetBlockHeader(height uint32) (*BlockHeader, error) {
	if req, err := http.NewRequest("GET", fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s/block/%d/header", w.Network, height), bytes.NewBuffer([]byte{})); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", w.ApiKey)
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == 404 {
				return nil, nil
			}
			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("failed to verify merkleroot for height %d because of an error: %v", height, resp.Status)
			}
			header := &BlockHeader{}
			if err := json.NewDecoder(resp.Body).Decode(header); err != nil {
				return nil, err
			}

			return header, nil
		}
	}
}

func (w *WhatsOnChain) IsValidRootForHeight(root *chainhash.Hash, height uint32) (bool, error) {
	if header, err := w.GetBlockHeader(height); err != nil {
		return false, err
	} else {
		return header.MerkleRoot.IsEqual(root), nil
	}
}
