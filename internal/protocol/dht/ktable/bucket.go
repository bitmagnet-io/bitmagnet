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
		btree:   btree.New(origin[:], k, true, true),
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
	var evictedID btree.NodeID
	it, ok := r.items[id]
	if ok {
		putResult = btree.PutAlreadyExists
	} else {
		putResult, evictedID = r.btree.Put(id[:])
		if putResult == btree.PutAlreadyExists {
			panic("shouldn't happen")
		}
	}
	if evictedID != nil {
		evictedItem, ok := r.items[protocol.MustNewIDFromByteSlice(evictedID)]
		if !ok {
			panic("shouldn't happen")
		}
		evictedItem.drop(errors.New("evicted"))
		delete(r.items, evictedItem.ID())
	}
	switch putResult {
	case btree.PutAccepted:
		it = r.newItem(id, input)
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

func (r bucketRoot[_, _, ItemPublic, _]) get(id ID) (it ItemPublic, _ bool) {
	if prv, ok := r.items[id]; ok {
		return prv.public(), true
	}
	return it, false
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
	if !r.btree.Drop(id[:]) {
		return false
	}
	if reason == nil {
		reason = ErrDropReasonNotProvided
	}
	r.items[id].drop(reason)
	delete(r.items, id)
	return true
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
