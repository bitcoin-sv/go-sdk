package overlay

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/util"
)

type Outpoint []byte

func NewOutpoint(txid []byte, vout uint32) *Outpoint {
	o := Outpoint(binary.LittleEndian.AppendUint32(util.ReverseBytes(txid), vout))
	return &o
}

func NewOutpointFromHash(txid *chainhash.Hash, vout uint32) *Outpoint {
	o := Outpoint(binary.LittleEndian.AppendUint32(txid.CloneBytes(), vout))
	return &o
}

func NewOutpointFromBytes(b []byte) *Outpoint {
	buf := make([]byte, 36)
	copy(buf, b)
	o := Outpoint(buf)
	return &o
}

func NewOutpointFromString(s string) (o *Outpoint, err error) {
	if len(s) < 66 {
		return nil, fmt.Errorf("invalid-string")
	}
	txid, err := hex.DecodeString(s[:64])
	if err != nil {
		return
	}
	vout, err := strconv.ParseUint(s[65:], 10, 32)
	if err != nil {
		return
	}
	origin := Outpoint(binary.LittleEndian.AppendUint32(util.ReverseBytes(txid), uint32(vout)))
	o = &origin
	return
}

func (o *Outpoint) String() string {
	return fmt.Sprintf("%x_%d", util.ReverseBytes((*o)[:32]), binary.LittleEndian.Uint32((*o)[32:]))
}

func (o *Outpoint) Txid() []byte {
	return util.ReverseBytes((*o)[:32])
}

func (o *Outpoint) TxidHash() *chainhash.Hash {
	hash, _ := chainhash.NewHash((*o)[:32])
	return hash
}

func (o *Outpoint) TxidHex() string {
	return hex.EncodeToString(o.Txid())
}

func (o *Outpoint) Vout() uint32 {
	return binary.LittleEndian.Uint32((*o)[32:])
}

func (o Outpoint) MarshalJSON() (bytes []byte, err error) {
	if len(o) != 36 {
		return []byte("null"), nil
	}
	return json.Marshal(o.String())
}

// UnmarshalJSON deserializes Origin to string
func (o *Outpoint) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	} else if op, err := NewOutpointFromString(x); err != nil {
		return err
	} else {
		*o = *op
		return nil
	}
}

func (o Outpoint) Value() (driver.Value, error) {
	b := make([]byte, 36)
	copy(b, o.Txid())
	binary.BigEndian.PutUint32(b[32:], o.Vout())
	return b, nil
}

func (o *Outpoint) Scan(value interface{}) error {
	if b, ok := value.([]byte); !ok || len(b) != 36 {
		return fmt.Errorf("invalid-outpoint")
	} else {
		outpoint := NewOutpoint(b[:32], binary.BigEndian.Uint32(b[32:]))
		*o = *outpoint
		return nil
	}
}
