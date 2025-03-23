package server

import (
	"encoding/binary"
	"sync"
)

type IDIssuer interface {
	Issue() string
}

type variantIDIssuer struct {
	mu   sync.Mutex
	buf  [binary.MaxVarintLen64]byte
	next uint64
}

func (i *variantIDIssuer) Issue() string {
	i.mu.Lock()
	n := binary.PutUvarint(i.buf[:], i.next)
	i.next++
	id := string(i.buf[:n])
	i.mu.Unlock()

	return id
}
