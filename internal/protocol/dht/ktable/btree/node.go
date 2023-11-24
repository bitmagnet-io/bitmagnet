package btree

type baseNode struct {
	origin NodeID
	bit    Bit
	path   Bits
}

func (n baseNode) N() int {
	return len(n.origin) * 8
}

type rootNode struct {
	baseNode
	k                int
	node             iNode
	splittingEnabled bool
	bucketCounts map[int]int
}

func New(origin NodeID, k int, splittingEnabled bool) Btree {
	return &rootNode{
		baseNode: baseNode{
			origin: origin,
		},
		node: emptyNode{
			baseNode: baseNode{
				origin: origin,
				path:   []Bit{},
			},
		},
		k:                k,
		splittingEnabled: splittingEnabled,
		bucketCounts: make(map[int]int, len(origin)*8),
	}
}

func (n *rootNode) Has(id NodeID) bool {
	return n.node.has(id.MustXor(n.origin))
}

func (n *rootNode) Put(id NodeID) PutResult {
	if id.Equals(n.origin) {
		return PutRejected
	}
	xor := id.MustXor(n.origin)
	if n.node.has(xor) {
		return PutAlreadyExists
	}
	// first check if we have a "bucket" with capacity
	if n.getBucketCount(xor) >= n.k {
		// then if "splitting" is enabled, check if we have fewer than k nodes closer to the origin
		if !n.splittingEnabled || n.node.countCloserThanSubpath(xor.Bits()) >= n.k {
			return PutRejected
		}
	}
	newNode, result := n.node.put(xor)
	if result == PutAccepted {
		n.node = newNode
		n.putBucketCount(xor)
	}
	return result
}

func (n *rootNode) putBucketCount(xor NodeID) {
	bucket := xor.Bits().LeadingZeros()
	if _, ok := n.bucketCounts[bucket]; !ok {
		n.bucketCounts[bucket] = 0
	}
	n.bucketCounts[bucket]++
}

func (n *rootNode) dropBucketCount(xor NodeID) {
	bucket := xor.Bits().LeadingZeros()
	n.bucketCounts[bucket]--
}

func (n *rootNode) getBucketCount(xor NodeID) int {
	bucket := xor.Bits().LeadingZeros()
	return n.bucketCounts[bucket]
}

func (n *rootNode) Drop(id NodeID) bool {
	xor := id.MustXor(n.origin)
	newNode, ok := n.node.drop(xor)
	if ok {
		n.node = newNode
		n.dropBucketCount(xor)
	}
	return ok
}

func (n *rootNode) furthest(thresholdXor NodeID) (NodeID, bool) {
	xor, ok := n.node.furthestXor()
	if ok {
		if xor.Bits().Cmp(thresholdXor.Bits()) <= 0 {
			return nil, false
		}
		return xor.MustXor(n.origin), true
	}
	return nil, false
}

func (n *rootNode) Count() int {
	return n.node.count()
}

func (n *rootNode) countCloserThan(id NodeID) int {
	return n.node.countCloserThanSubpath(id.Bits())
}

func (n *rootNode) Closest(id NodeID, count int) []NodeID {
	xors := n.node.xorsClosestToSubpath(id.MustXor(n.origin).Bits(), count)
	ids := make([]NodeID, len(xors))
	for i, xor := range xors {
		ids[i] = xor.MustXor(n.origin)
	}
	return ids
}

type iNode interface {
	has(NodeID) bool
	put(NodeID) (iNode, PutResult)
	drop(NodeID) (iNode, bool)
	furthestXor() (NodeID, bool)
	allXors() []NodeID
	any() bool
	count() int
	countCloserThanSubpath(Bits) int
	countAtSubpath(Bits) int
	xorsClosestToSubpath(Bits, int) []NodeID
}

type emptyNode struct {
	baseNode
}

func (n emptyNode) has(NodeID) bool {
	return false
}

func (n emptyNode) put(xor NodeID) (iNode, PutResult) {
	return leafNode{
		baseNode: n.baseNode,
		xor:      xor,
	}, PutAccepted
}

func (n emptyNode) drop(NodeID) (iNode, bool) {
	return n, false
}

func (n emptyNode) furthestXor() (NodeID, bool) {
	return nil, false
}

func (n emptyNode) allXors() []NodeID {
	return nil
}

func (n emptyNode) any() bool {
	return false
}

func (n emptyNode) count() int {
	return 0
}

func (n emptyNode) countCloserThanSubpath(Bits) int {
	return 0
}

func (n emptyNode) countAtSubpath(Bits) int {
	return 0
}

func (n emptyNode) xorsClosestToSubpath(Bits, int) []NodeID {
	return nil
}

type leafNode struct {
	baseNode
	xor NodeID
}

func (n leafNode) has(xor NodeID) bool {
	return n.xor.Equals(xor)
}

func (n leafNode) put(xor NodeID) (iNode, PutResult) {
	if n.xor.Equals(xor) {
		return n, PutAlreadyExists
	}
	// unlike the emptyNode and leafNode, the branchNode is initialized as a pointer as it may have many thousands of descendents
	initNode := &branchNode{
		baseNode: n.baseNode,
		branches: map[Bit]iNode{
			Bit0: emptyNode{
				baseNode: baseNode{
					origin: n.origin,
					bit:    Bit0,
					path:   appendToPath(n.path, Bit0),
				},
			},
			Bit1: emptyNode{
				baseNode: baseNode{
					origin: n.origin,
					bit:    Bit1,
					path:   appendToPath(n.path, Bit1),
				},
			},
		},
		counts: map[Bit]int{
			false: 0,
			true:  0,
		},
	}
	newNode, existingPutResult := initNode.put(n.xor)
	if existingPutResult != PutAccepted {
		panic("unexpected PutResult")
	}
	return newNode.put(xor)
}

func (n leafNode) drop(xor NodeID) (iNode, bool) {
	if n.xor.Equals(xor) {
		return emptyNode{
			baseNode: n.baseNode,
		}, true
	}
	return n, false
}

func (n leafNode) furthestXor() (NodeID, bool) {
	return n.xor, true
}

func (n leafNode) allXors() []NodeID {
	return []NodeID{n.xor}
}

func (n leafNode) any() bool {
	return true
}

func (n leafNode) count() int {
	return 1
}

func (n leafNode) countCloserThanSubpath(path Bits) int {
	for i, bit := range path {
		if bit && !n.xor.GetBit(i+len(n.path)) {
			return 0
		}
	}
	return 1
}

func (n leafNode) countAtSubpath(path Bits) int {
	for i, bit := range path {
		if n.xor.GetBit(i+len(n.path)) != bit {
			return 0
		}
	}
	return 1
}

func (n leafNode) xorsClosestToSubpath(_ Bits, count int) []NodeID {
	if count < 1 {
		return nil
	}
	return []NodeID{n.xor}
}

type branchNode struct {
	baseNode
	branches map[Bit]iNode
	counts   map[Bit]int
}

func (n branchNode) has(xor NodeID) bool {
	return n.branches[xor.GetBit(len(n.path))].has(xor)
}

func (n branchNode) put(xor NodeID) (iNode, PutResult) {
	bit := xor.GetBit(len(n.path))
	branch := n.branches[bit]
	newBranch, result := branch.put(xor)
	if result == PutAccepted {
		n.branches[bit] = newBranch
		n.counts[bit] = newBranch.count()
	}
	return n, result
}

func (n branchNode) drop(xor NodeID) (iNode, bool) {
	bit := xor.GetBit(len(n.path))
	branch := n.branches[bit]
	newBranch, ok := branch.drop(xor)
	if ok {
		newCount := newBranch.count()
		if newCount == 0 {
			if n.counts[!bit] == 0 {
				return emptyNode{baseNode: n.baseNode}, true
			}
			if n.counts[!bit] == 1 {
				allXors := n.branches[!bit].allXors()
				if len(allXors) != 1 {
					panic("unexpected condition")
				}
				return leafNode{
					baseNode: n.baseNode,
					xor:      allXors[0],
				}, true
			}
		}
		n.branches[bit] = newBranch
		n.counts[bit] = newCount
	}
	return n, ok
}

func (n branchNode) furthestXor() (NodeID, bool) {
	xor, ok := n.branches[Bit1].furthestXor()
	if ok {
		return xor, true
	}
	return n.branches[Bit0].furthestXor()
}

func (n branchNode) allXors() []NodeID {
	var xors []NodeID
	xors = append(xors, n.branches[Bit0].allXors()...)
	xors = append(xors, n.branches[Bit1].allXors()...)
	return xors
}

func (n branchNode) any() bool {
	return n.branches[Bit0].any() || n.branches[Bit1].any()
}

func (n branchNode) count() int {
	return n.counts[Bit0] + n.counts[Bit1]
}

func (n branchNode) countCloserThanSubpath(path Bits) int {
	if len(path) == 0 {
		return 0
	}
	if len(path) == 1 {
		switch path[0] {
		case Bit0:
			return n.counts[Bit0]
		case Bit1:
			return n.counts[Bit0] + n.counts[Bit1]
		}
	}
	if path[0] == Bit0 {
		return n.branches[Bit0].countCloserThanSubpath(path[1:])
	}
	return n.counts[Bit0] + n.branches[Bit1].countCloserThanSubpath(path[1:])
}

func (n branchNode) countAtSubpath(path Bits) int {
	if len(path) == 0 {
		return n.count()
	}
	return n.branches[path[0]].countAtSubpath(path[1:])
}

func (n branchNode) xorsClosestToSubpath(path Bits, count int) []NodeID {
	if len(path) == 0 {
		path = appendToPath(path, Bit0)
	}
	closest := n.branches[path[0]].xorsClosestToSubpath(path[1:], count)
	if len(closest) < count {
		closest = append(closest, n.branches[!path[0]].xorsClosestToSubpath(nil, count-len(closest))...)
	}
	return closest
}
