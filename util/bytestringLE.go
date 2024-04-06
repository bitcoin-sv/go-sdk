package util

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// ByteStringLE is a byte array that serializes to hex
type ByteStringLE []byte

func NewByteStringLEFromHex(s string) ByteStringLE {
	b, _ := hex.DecodeString(s)
	return ByteStringLE(b)
}

// MarshalJSON serializes ByteArray to hex
func (s ByteStringLE) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON deserializes ByteArray to hex
func (s *ByteStringLE) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if len(x) > 0 {
		str, err := hex.DecodeString(x)
		if err != nil {
			return err
		}
		str = ReverseBytes(str)
		*s = ByteStringLE(str)
	} else {
		*s = nil
	}
	return nil
}

func (s *ByteStringLE) String() string {
	return hex.EncodeToString(ReverseBytes(*s))
}

func (s ByteStringLE) Value() (driver.Value, error) {
	return []byte(s), nil
}

func (s *ByteStringLE) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*s = b
	return nil
}
