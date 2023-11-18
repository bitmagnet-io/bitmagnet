package ktable

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"
)

type keyspace[
	Input any,
	Option any,
	ItemPublic keyspaceItem,
	ItemPrivate keyspaceItemPrivate[Input, Option, ItemPublic],
] struct {
	origin  ID
	k       int
	btree   btree.Btree
	items   map[ID]ItemPrivate
	factory func(ID, Input) ItemPrivate
}

func newKeyspace[
	Input any,
	Option any,
	ItemPublic keyspaceItem,
	ItemPrivate keyspaceItemPrivate[Input, Option, ItemPublic],
](
	origin ID,
	k int,
	newItem func(ID, Input) ItemPrivate,
) keyspace[Input, Option, ItemPublic, ItemPrivate] {
	return keyspace[Input, Option, ItemPublic, ItemPrivate]{
		origin:  origin,
		k:       k,
		btree:   btree.New(origin[:], k, true),
		items:   make(map[ID]ItemPrivate, k),
		factory: newItem,
	}
}

func (k keyspace[Input, _, _, ItemPrivate]) newItem(id ID, input Input) ItemPrivate {
	return k.factory(id, input)
}

func (k keyspace[_, _, _, _]) count() int {
	return len(k.items)
}

func (k keyspace[_, _, ItemPublic, _]) get(id ID) (it ItemPublic, ok bool) {
	prv, ok := k.items[id]
	if ok {
		it = prv.public()
	}
	return
}

func (k keyspace[Input, Option, _, ItemPrivate]) put(id ID, input Input, options ...Option) btree.PutResult {
	var it ItemPrivate
	var putResult btree.PutResult
	it, ok := k.items[id]
	putResult = k.btree.Put(id[:])
	switch putResult {
	case btree.PutAccepted:
		if !ok {
			it = k.newItem(id, input)
		}
		k.items[id] = it
	case btree.PutAlreadyExists:
		it.update(input)
	default:
		return putResult
	}
	for _, o := range options {
		it.apply(o)
	}
	return putResult
}

func (k keyspace[_, _, ItemPublic, _]) getRandom(n int) []ItemPublic {
	if n > len(k.items) {
		n = len(k.items)
	}
	items := make([]ItemPublic, 0, n)
	for _, prv := range k.items {
		if len(items) >= n {
			break
		}
		items = append(items, prv.public())
	}
	return items
}

var ErrDropReasonNotProvided = errors.New("drop reason not provided")

func (k keyspace[_, _, _, _]) drop(id ID, reason error) bool {
	it, ok := k.items[id]
	if !ok {
		return false
	}
	if reason == nil {
		reason = ErrDropReasonNotProvided
	}
	it.drop(reason)
	delete(k.items, id)
	return k.btree.Drop(id[:])
}

func (k keyspace[_, _, ItemPublic, _]) getClosest(id ID) []ItemPublic {
	if it, ok := k.items[id]; ok {
		return []ItemPublic{it.public()}
	}
	closestIDs := k.btree.Closest(id[:], 8)
	closest := make([]ItemPublic, 0, len(closestIDs))
	for _, id := range closestIDs {
		closest = append(closest, k.items[protocol.MustNewIDFromByteSlice(id)].public())
	}
	return closest
}

type keyspaceItem interface {
	ID() ID
}

type keyspaceItemPrivate[
	Input any,
	Option any,
	Public keyspaceItem,
] interface {
	keyspaceItem
	update(Input)
	apply(Option)
	drop(error)
	public() Public
}
