package protocol

import (
	crand "crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	"io"
	"strings"
)

// There are 2 main conventions for encodeing client and client version information into the client ID,
// Azureus-style and Shadow's-style.
//
// Azureus-style uses the following encoding: '-', two characters for client id, four ascii digits for version
// number, '-', followed by random numbers.
//
// For example: '-AZ2060-'...
//
// https://wiki.theory.org/BitTorrentSpecification#peer_id
//
// We encode the version number as:
// - First two digits for the major version number
// - Last two digits for the minor version number
// - Patch version number is not encoded.
const idClientPart = "-BM0001-"

func RandomNodeID() (id ID) {
	_, _ = crand.Read(id[:])
	return
}

// RandomNodeIDWithClientSuffix generates a node ID for the DHT client.
// We use a random byte string with the client ID encoded at the end, to allow identifying other bitmagnet instances in the wild.
// A suffix is used instead of a prefix, which would be incompatible with DHT, where ID prefixes are used for computing the distance metric).
func RandomNodeIDWithClientSuffix() (id ID) {
	_, _ = crand.Read(id[:])
	for i := 0; i < len(idClientPart); i++ {
		id[20-len(idClientPart)+i] = idClientPart[i]
	}
	return
}

type ID [20]byte

func ParseID(str string) (ID, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(str, "0x"))
	if err != nil {
		return ID{}, err
	}
	if len(b) != 20 {
		return ID{}, errors.New("hash string must be 20 bytes")
	}
	var id ID
	copy(id[:], b)
	return id, nil
}

func MustParseID(str string) ID {
	id, err := ParseID(str)
	if err != nil {
		panic(err)
	}
	return id
}

func NewIDFromRawString(s string) (id ID) {
	if n := copy(id[:], s); n != 20 {
		panic(n)
	}
	return
}

func NewIDFromByteSlice(b []byte) (id ID, _ error) {
	if n := copy(id[:], b); n != 20 {
		return id, errors.New("must be 20 bytes")
	}
	return
}

func MustNewIDFromByteSlice(b []byte) ID {
	id, err := NewIDFromByteSlice(b)
	if err != nil {
		panic(err)
	}
	return id
}

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id ID) Int160() Int160 {
	return NewInt160FromByteArray(id)
}

func (id ID) IsZero() bool {
	return id == [20]byte{}
}

func (id ID) GetBit(i int) bool {
	return id[i/8]>>(7-uint(i%8))&1 == 1
}

func (id ID) Bytes() []byte {
	return id[:]
}

func (b *ID) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid bytes type")
	}
	copy(b[:], v)
	return nil
}

func (b ID) Value() (driver.Value, error) {
	return b[:], nil
}

func (b ID) MarshalBinary() ([]byte, error) {
	return b[:], nil
}

func (b *ID) UnmarshalBinary(data []byte) error {
	if len(data) != 20 {
		return errors.New("invalid ID length")
	}
	copy(b[:], data)
	return nil
}

func (id ID) MarshalBencode() ([]byte, error) {
	return []byte("20:" + string(id[:])), nil
}

func (id *ID) UnmarshalBencode(b []byte) error {
	var s string
	if err := bencode.Unmarshal(b, &s); err != nil {
		return err
	}
	if n := copy(id[:], s); n != 20 {
		return fmt.Errorf("string has wrong length: %d", n)
	}
	return nil
}

func (b ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tb, err := ParseID(s)
	if err != nil {
		return err
	}
	*b = tb
	return nil
}

func (b *ID) UnmarshalGQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		tb, err := ParseID(input)
		if err != nil {
			return err
		}
		*b = tb
		return nil
	default:
		return errors.New("invalid hash type")
	}
}

func (b ID) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + b.String() + `"`))
}

type MutableID ID

func (id *MutableID) SetBit(i int, v bool) {
	if v {
		id[i/8] |= 1 << (7 - uint(i%8))
	} else {
		id[i/8] &= ^(1 << (7 - uint(i%8)))
	}
}

func RandomPeerID() ID {
	clientID := RandomNodeID()
	i := 0
	for _, c := range idClientPart {
		clientID[i] = byte(c)
		i++
	}
	return clientID
}
