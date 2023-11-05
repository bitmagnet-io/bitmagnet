package ktable

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"
	"math/rand"
)

type bucketRoot[
	Input any,
	Option any,
	ItemPublic bucketItem,
	ItemPrivate bucketItemPrivate[Input, Option, ItemPublic],
] struct {
	origin  ID
	k       int
	btree   btree.Btree
	items   map[ID]ItemPrivate
	newItem func(ID, Input) ItemPrivate
}

func newBucketRoot[
	Input any,
	Option any,
	ItemPublic bucketItem,
	ItemPrivate bucketItemPrivate[Input, Option, ItemPublic],
](
	origin ID,
	k int,
	newItem func(ID, Input) ItemPrivate,
) bucketRoot[Input, Option, ItemPublic, ItemPrivate] {
	return bucketRoot[Input, Option, ItemPublic, ItemPrivate]{
		origin:  origin,
		k:       k,
		btree:   btree.New(origin[:], k, true),
		items:   make(map[ID]ItemPrivate, k),
		newItem: newItem,
	}
}

func (r bucketRoot[_, _, _, _]) count() int {
	return len(r.items)
}

func (r bucketRoot[Input, Option, _, ItemPrivate]) put(id ID, input Input, options ...Option) btree.PutResult {
	var it ItemPrivate
	var putResult btree.PutResult
	it, ok := r.items[id]
	putResult = r.btree.Put(id[:])
	switch putResult {
	case btree.PutAccepted:
		if !ok {
			it = r.newItem(id, input)
		}
		r.items[id] = it
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

func (r bucketRoot[_, _, ItemPublic, _]) getRandom(n int) []ItemPublic {
	if n > len(r.items) {
		n = len(r.items)
	}
	items := make([]ItemPublic, 0, n)
	for _, prv := range r.items {
		if len(items) >= n {
			break
		}
		items = append(items, prv.public())
	}
	return items
}

var ErrDropReasonNotProvided = errors.New("drop reason not provided")

func (r bucketRoot[_, _, _, _]) drop(id ID, reason error) bool {
	it, ok := r.items[id]
	if !ok {
		return false
	}
	if reason == nil {
		reason = ErrDropReasonNotProvided
	}
	it.drop(reason)
	delete(r.items, id)
	return r.btree.Drop(id[:])
}

func (r bucketRoot[_, _, ItemPublic, _]) getClosest(id ID) []ItemPublic {
	if it, ok := r.items[id]; ok {
		return []ItemPublic{it.public()}
	}
	closestIDs := r.btree.Closest(id[:], 8)
	closest := make([]ItemPublic, 0, len(closestIDs))
	for _, id := range closestIDs {
		closest = append(closest, r.items[protocol.MustNewIDFromByteSlice(id)].public())
	}
	return closest
}

func (r bucketRoot[_, _, _, _]) generateRandomID() ID {
	lengths := r.btree.EmptiestPrefixLengths()
	length := lengths[rand.Intn(rand.Intn(len(lengths))+1)]
	id := protocol.MutableID(protocol.RandomNodeID())
	for i := 0; i < length; i++ {
		id.SetBit(i, r.origin.GetBit(i))
	}
	if length < 160 {
		id.SetBit(length, !r.origin.GetBit(length))
	}
	return ID(id)
}

type bucketItem interface {
	ID() ID
}

type bucketItemPrivate[
	Input any,
	Option any,
	Public bucketItem,
] interface {
	bucketItem
	update(Input)
	apply(Option)
	drop(error)
	public() Public
}
