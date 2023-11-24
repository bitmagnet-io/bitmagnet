package server

import (
	"encoding/binary"
	"sync"
)

type IdIssuer interface {
	Issue() string
}

type variantIdIssuer struct {
	mu   sync.Mutex
	buf  [binary.MaxVarintLen64]byte
	next uint64
}

func (i *variantIdIssuer) Issue() string {
	i.mu.Lock()
	n := binary.PutUvarint(i.buf[:], i.next)
	i.next++
	id := string(i.buf[:n])
	i.mu.Unlock()
	return id
}
