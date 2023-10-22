package ktable

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"
	"net/netip"
)

type Command interface {
	exec(*table)
}

type CommandReturn[T any] interface {
	Command
	Query[T]
}

type command[T any] struct {
	fn func(*table) T
}

func (c command[T]) exec(t *table) {
	c.fn(t)
}

func (c command[T]) execReturn(t *table) T {
	return c.fn(t)
}

var _ CommandReturn[btree.PutResult] = PutPeer{}

type PutPeer struct {
	ID      ID
	Addr    netip.AddrPort
	Options []PeerOption
}

func (c PutPeer) execReturn(t *table) btree.PutResult {
	if !c.Addr.IsValid() {
		return btree.PutRejected
	}
	return t.peers.put(c.ID, c.Addr, c.Options...)
}

func (c PutPeer) exec(t *table) {
	c.execReturn(t)
}

var _ CommandReturn[bool] = DropPeer{}

type DropPeer struct {
	ID     ID
	Reason error
}

func (c DropPeer) execReturn(t *table) bool {
	return t.peers.drop(c.ID, c.Reason)
}

func (c DropPeer) exec(t *table) {
	c.execReturn(t)
}

var _ CommandReturn[bool] = DropAddr{}

type DropAddr struct {
	Addr   netip.Addr
	Reason error
}

func (c DropAddr) execReturn(t *table) bool {
	id, ok := t.addrs.getPeerIDForAddr(c.Addr)
	if !ok {
		return false
	}
	return t.peers.drop(id, c.Reason)
}

func (c DropAddr) exec(t *table) {
	c.execReturn(t)
}

var _ CommandReturn[btree.PutResult] = PutHash{}

type PutHash struct {
	ID      ID
	Peers   []HashPeer
	Options []HashOption
}

func (c PutHash) execReturn(t *table) btree.PutResult {
	return t.hashes.put(c.ID, c.Peers, c.Options...)
}

func (c PutHash) exec(t *table) {
	c.execReturn(t)
}
