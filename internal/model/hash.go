package model

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type Hash20 [20]byte

func NewHash20FromString(s string) (h Hash20, err error) {
	b, decodeErr := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if decodeErr != nil {
		err = decodeErr
		return
	}
	if len(b) != 20 {
		err = errors.New("hash string must be 20 bytes")
		return
	}
	var tb Hash20
	copy(tb[:], b)
	return tb, nil
}

func (b *Hash20) String() string {
	return hex.EncodeToString(b[:])
}

func (b *Hash20) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid bytes type")
	}
	copy(b[:], v)
	return nil
}

func (b Hash20) Value() (driver.Value, error) {
	return b[:], nil
}

func (b *Hash20) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *Hash20) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tb, err := NewHash20FromString(s)
	if err != nil {
		return err
	}
	*b = tb
	return nil
}

func (b *Hash20) UnmarshalGQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		tb, err := NewHash20FromString(input)
		if err != nil {
			return err
		}
		*b = tb
		return nil
	default:
		return errors.New("invalid hash type")
	}
}

func (b Hash20) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + b.String() + `"`))
}
