package broadcaster

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

type MockFailureClient struct{}

func (m *MockFailureClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
	}, nil
}

type MockSuccessClient struct{}

func (m *MockSuccessClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"txid":"4d76b00f29e480e0a933cef9d9ffe303d6ab919e2cdb265dd2cea41089baa85a"}`)),
	}, nil
}

type MockNetworkErrorClient struct{}

func (m *MockNetworkErrorClient) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("network error")
}

type MockBadRequestClient struct{}

func (m *MockBadRequestClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader("Bad Request")),
	}, nil
}

type MockUnauthorizedClient struct{}

func (m *MockUnauthorizedClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 401,
		Body:       io.NopCloser(strings.NewReader("Unauthorized")),
	}, nil
}

type MockBodyReadErrorClient struct{}

type ErrorReader struct{}

func (e *ErrorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}

func (e *ErrorReader) Close() error {
	return nil
}

func (m *MockBodyReadErrorClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       &ErrorReader{},
	}, nil
}

func TestWhatsOnChainBroadcast(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
		// For simplicity, we'll use an empty transaction
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockSuccessClient{},
	}

	success, failure := b.Broadcast(tx)
	require.NotNil(t, success)
	require.Nil(t, failure)
	require.Equal(t, "f702453dd03b0f055e5437d76128141803984fb10acb85fc3b2184fae2f3fa78", success.Txid)
}

func TestWhatsOnChainBroadcastFailure(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockFailureClient{},
	}

	success, failure := b.Broadcast(tx)
	require.Nil(t, success)
	require.NotNil(t, failure)
	require.Equal(t, "500", failure.Code)
	require.Equal(t, "Internal Server Error", failure.Description)
}

func TestWhatsOnChainBroadcastClientError(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockNetworkErrorClient{},
	}

	success, failure := b.Broadcast(tx)
	require.Nil(t, success)
	require.NotNil(t, failure)
	require.Contains(t, failure.Description, "network error")
}

func TestWhatsOnChainBroadcastBadRequest(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockBadRequestClient{},
	}

	success, failure := b.Broadcast(tx)
	require.Nil(t, success)
	require.NotNil(t, failure)
	require.Equal(t, "400", failure.Code)
	require.Equal(t, "Bad Request", failure.Description)
}

func TestWhatsOnChainBroadcastUnauthorized(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "invalid_api_key",
		Client:  &MockUnauthorizedClient{},
	}

	success, failure := b.Broadcast(tx)
	require.Nil(t, success)
	require.NotNil(t, failure)
	require.Equal(t, "401", failure.Code)
	require.Equal(t, "Unauthorized", failure.Description)
}

func TestWhatsOnChainBroadcastBodyReadError(t *testing.T) {
	tx := &transaction.Transaction{
		// Populate with valid data
	}

	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockBodyReadErrorClient{},
	}

	success, failure := b.Broadcast(tx)
	require.Nil(t, success)
	require.NotNil(t, failure)
	require.Equal(t, "500", failure.Code)
	require.Equal(t, "unknown error", failure.Description)
}
