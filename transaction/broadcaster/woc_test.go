package broadcaster

import (
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
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

type MockSuccessClient struct{}

func (m *MockSuccessClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"txid":"4d76b00f29e480e0a933cef9d9ffe303d6ab919e2cdb265dd2cea41089baa85a"}`)),
	}, nil
}

func TestWhatsOnChainBroadcastFail(t *testing.T) {
	tx, err := transaction.NewTransactionFromHex("0100000001a9b0c5a2437042e5d0c6288fad6abc2ef8725adb6fef5f1bab21b2124cfb7cf6dc9300006a47304402204c3f88aadc90a3f29669bba5c4369a2eebc10439e857a14e169d19626243ffd802205443013b187a5c7f23e2d5dd82bc4ea9a79d138a3dc6cae6e6ef68874bd23a42412103fd290068ae945c23a06775de8422ceb6010aaebab40b78e01a0af3f1322fa861ffffffff010000000000000000b1006a0963657274696861736822314c6d763150594d70387339594a556e374d3948565473446b64626155386b514e4a4032356163343531383766613035616532626436346562323632386666336432666636646338313665383335376364616366343765663862396331656433663531403064383963343363343636303262643865313831376530393137313736343134353938373337623161663865363939343930646364653462343937656338643300000000")
	require.NoError(t, err)
	b := &WhatsOnChain{
		Network: "main",
		ApiKey:  "",
		Client:  &MockFailureClient{},
	}

	success, failure := tx.Broadcast(b)
	require.Nil(t, success)
	require.NotNil(t, failure)

	b.Client = &MockSuccessClient{}

	success, failure = tx.Broadcast(b)
	require.NotNil(t, success)
	require.Nil(t, failure)
	require.Equal(t, "4d76b00f29e480e0a933cef9d9ffe303d6ab919e2cdb265dd2cea41089baa85a", success.Txid)
}
