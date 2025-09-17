package api_key

import (
	"encoding/base64"
	"encoding/binary"
)

type keyData struct {
	id     int
	secret []byte
}

func (k keyData) encode() string {
	bytes := make([]byte, 0, 8+secretLength)
	binary.BigEndian.AppendUint64(bytes, uint64(k.id))
	bytes = append(bytes, k.secret...)

	return base64.StdEncoding.EncodeToString(bytes)
}

func (k *keyData) decode(key string) error {
	rawBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	if len(rawBytes) != 8+secretLength {
		return ErrLength
	}

	k.id = int(binary.BigEndian.Uint64(rawBytes[:8]))
	k.secret = rawBytes[8:]

	return nil
}
