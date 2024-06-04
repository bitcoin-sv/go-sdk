package script_test

import (
	"fmt"
	"testing"

	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/stretchr/testify/assert"
)

func TestEncodeBIP276(t *testing.T) {
	t.Parallel()

	t.Run("valid encode (mainnet)", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixScript,
				Version: script.CurrentVersion,
				Network: script.NetworkMainnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "bitcoin-script:010166616b65207363726970746f0cd86a", s)
	})

	t.Run("valid encode (testnet)", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixScript,
				Version: script.CurrentVersion,
				Network: script.NetworkTestnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "bitcoin-script:020166616b65207363726970742577a444", s)
	})

	t.Run("invalid version = 0", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixScript,
				Version: 0,
				Network: script.NetworkMainnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "ERROR", s)
	})

	t.Run("invalid version > 255", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixScript,
				Version: 256,
				Network: script.NetworkMainnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "ERROR", s)
	})

	t.Run("invalid network = 0", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixScript,
				Version: script.CurrentVersion,
				Network: 0,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "ERROR", s)
	})

	t.Run("different prefix", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  "different-prefix",
				Version: script.CurrentVersion,
				Network: script.NetworkMainnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "different-prefix:010166616b6520736372697074effdb090", s)
	})

	t.Run("template prefix", func(t *testing.T) {
		s := script.EncodeBIP276(
			script.BIP276{
				Prefix:  script.PrefixTemplate,
				Version: script.CurrentVersion,
				Network: script.NetworkMainnet,
				Data:    []byte("fake script"),
			},
		)

		assert.Equal(t, "bitcoin-template:010166616b65207363726970749e31aa72", s)
	})
}

func TestDecodeBIP276(t *testing.T) {
	t.Parallel()

	t.Run("valid decode", func(t *testing.T) {
		script, err := script.DecodeBIP276("bitcoin-script:010166616b65207363726970746f0cd86a")
		assert.NoError(t, err)
		assert.Equal(t, `"bitcoin-script"`, fmt.Sprintf("%q", script.Prefix))
		assert.Equal(t, 1, script.Network)
		assert.Equal(t, 1, script.Version)
		assert.Equal(t, "fake script", string(script.Data))
	})

	t.Run("invalid decode", func(t *testing.T) {
		script, err := script.DecodeBIP276("bitcoin-script:01")
		assert.Error(t, err)
		assert.Nil(t, script)
	})

	t.Run("valid format, bad checksum", func(t *testing.T) {
		script, err := script.DecodeBIP276("bitcoin-script:010166616b65207363726970746f0cd8")
		assert.Error(t, err)
		assert.Nil(t, script)
	})
}
