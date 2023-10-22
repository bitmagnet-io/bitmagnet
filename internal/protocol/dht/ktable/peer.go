package ktable

import (
	"net/netip"
	"time"
)

type peerBucketRoot struct {
	bucketRoot[netip.AddrPort, PeerOption, Peer, *peer]
}

func (b *peerBucketRoot) getLastRespondedBefore(t time.Time) []Peer {
	var peers []Peer
	for _, p := range b.items {
		if p.lastRespondedAt.Before(t) {
			peers = append(peers, p)
		}
	}
	return peers
}

func (b *peerBucketRoot) getCandidatesForSampleInfoHashes(n int) []*peer {
	var candidates []*peer
	for _, p := range b.items {
		if !p.IsSampleInfoHashesCandidate() {
			continue
		}
		candidates = append(candidates, p)
		if len(candidates) == n {
			break
		}
	}
	return candidates
}

type Peer interface {
	bucketItem
	Addr() netip.AddrPort
	Time() time.Time
	Dropped() bool
	IsSampleInfoHashesCandidate() bool
}

type peer struct {
	id              ID
	addr            netip.AddrPort
	discoveredAt    time.Time
	lastRespondedAt time.Time
	dropReason      error
	peerBep51Info
	reverseMap *reverseMap
}

type peerBep51Info struct {
	support                  protocolSupport
	sampledNum               int
	lastSampledNum           int
	totalNum                 int
	nextSampleInfoHashesTime time.Time
}

// todo: we might want some time threshold for rechecking the total?
func (i *peerBep51Info) noneRemaining() bool {
	return i.support == protocolSupportYes && (i.totalNum-i.sampledNum < 1)
}

type PeerOption interface {
	apply(*peer)
}

type peerOption struct {
	fn func(*peer)
}

func (p peerOption) apply(peer *peer) {
	p.fn(peer)
}

type protocolSupport int

const (
	protocolSupportUnknown protocolSupport = iota
	protocolSupportYes
	protocolSupportNo
)

var _ bucketItemPrivate[netip.AddrPort, PeerOption, Peer] = (*peer)(nil)

func (p *peer) update(addr netip.AddrPort) {
	if p.addr != addr {
		p.reverseMap.dropAddr(p.addr.Addr())
	}
	p.addr = addr
	p.reverseMap.putAddrPeerID(addr.Addr(), p.id)
}

func (p *peer) apply(option PeerOption) {
	option.apply(p)
}

func (p *peer) drop(reason error) {
	p.dropReason = reason
	p.reverseMap.dropAddr(p.addr.Addr())
}

func (p *peer) public() Peer {
	return p
}

func (p *peer) ID() ID {
	return p.id
}

func (p *peer) Addr() netip.AddrPort {
	return p.addr
}

func (p *peer) Time() time.Time {
	return p.lastRespondedAt
}

func (p *peer) Dropped() bool {
	return p.dropReason != nil
}

func (p *peer) IsSampleInfoHashesCandidate() bool {
	now := time.Now()
	threshold := now.Add(-(5 * time.Second))
	return p.peerBep51Info.support != protocolSupportNo &&
		p.peerBep51Info.nextSampleInfoHashesTime.Before(now) &&
		p.lastRespondedAt.Before(threshold)
}

func (p *peer) olderThan(t time.Time) bool {
	return p.lastRespondedAt.Before(t)
}

func PeerResponded() PeerOption {
	return peerOption{
		fn: func(p *peer) {
			p.lastRespondedAt = time.Now()
		},
	}
}

func PeerBep51Support(supported bool) PeerOption {
	return peerOption{
		fn: func(p *peer) {
			s := protocolSupportNo
			if supported {
				s = protocolSupportYes
			}
			p.peerBep51Info.support = s
		},
	}
}

func PeerSampleInfoHashesRes(sampledNum int, totalNum int, nextSampleTime time.Time) PeerOption {
	return peerOption{
		fn: func(p *peer) {
			p.peerBep51Info.sampledNum += sampledNum
			p.peerBep51Info.lastSampledNum = sampledNum
			p.peerBep51Info.totalNum += totalNum
			// a crude way of deprioritizing peers that gave us no samples:
			if sampledNum == 0 {
				now := time.Now()
				if nextSampleTime.Before(now) {
					nextSampleTime = now
				}
				nextSampleTime = nextSampleTime.Add(5 * time.Minute)
			}
			p.peerBep51Info.nextSampleInfoHashesTime = nextSampleTime
		},
	}
}
