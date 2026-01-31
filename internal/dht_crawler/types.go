package dht_crawler

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Runner runner.Provider

type nodeHasPeersForHash struct {
	infoHash               protocol.ID
	node                   netip.AddrPort
	isVerifiedAbsentFromDB bool
}

type infoHashWithMetaInfo struct {
	nodeHasPeersForHash
	metaInfo metainfo.Info
}

type infoHashWithPeers struct {
	nodeHasPeersForHash
	peers []netip.AddrPort
}
