package api_key

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	secretLength  = 12
	keyLength     = 22
	decodedLength = (keyLength*5 + 7) / 8
)

type Secret struct {
	Secret []byte
	Hash   []byte
}

func NewSecret() Secret {
	bytes := make([]byte, secretLength)
	_, _ = rand.Read(bytes)

	hash, _ := bcrypt.GenerateFromPassword(
		bytes,
		bcrypt.DefaultCost,
	)

	return Secret{
		Secret: bytes,
		Hash:   hash,
	}
}

type KeyData struct {
	ID     int
	Secret []byte
}

func (k KeyData) Encode() string {
	bytes := make([]byte, 0, secretLength+4)
	bytes = append(bytes, k.Secret...)
	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, uint32(k.ID))
	bytes = append(bytes, idBytes...)

	return base62Encode(bytes)
}

func (k *KeyData) Decode(key string) error {
	bytes, err := base62Decode(key)
	if err != nil {
		return err
	}

	if len(bytes) != secretLength+4 {
		return ErrDecode
	}

	k.Secret = bytes[:secretLength]
	k.ID = int(binary.LittleEndian.Uint32(bytes[secretLength:]))

	return nil
}

var maxLen = len(new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(secretLength+4)), nil).Text(62))

func base62Encode(data []byte) string {
	bi := new(big.Int).SetBytes(data)
	// Pad with leading zeros
	s := bi.Text(62)
	if len(s) < maxLen {
		s = strings.Repeat("0", maxLen-len(s)) + s
	}

	return s
}

func base62Decode(s string) ([]byte, error) {
	if len(s) != keyLength {
		return nil, ErrDecode
	}

	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return nil, ErrDecode
		}
	}

	bi, ok := new(big.Int).SetString(s, 62)
	if !ok {
		return nil, ErrDecode
	}

	b := bi.Bytes()
	if len(b) < decodedLength {
		// Pad with leading zeros
		padded := make([]byte, decodedLength)
		copy(padded[decodedLength-len(b):], b)

		return padded, nil
	}

	return b, nil
}
