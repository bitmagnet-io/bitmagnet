package dht

import "github.com/anacrolix/dht/v2/krpc"

const (
	QPing             = "ping"
	QFindNode         = "find_node"
	QGetPeers         = "get_peers"
	QAnnouncePeer     = "announce_peer"
	QSampleInfohashes = "sample_infohashes"
)

type RecvMsg struct {
	Msg  krpc.Msg
	From krpc.NodeAddr
}

// AnnouncePort returns the torrent port for the message.
// There is an optional argument called implied_port which value is either 0 or 1.
// If it is present and non-zero, the port argument should be ignored and the source port of the UDP packet should be used as the peer's port instead.
// This is useful for peers behind a NAT that may not know their external port, and supporting uTP,
// they accept incoming connections on the same port as the DHT port.
// https://www.bittorrent.org/beps/bep_0005.html
func (m RecvMsg) AnnouncePort() int {
	port := m.From.Port
	args := m.Msg.A
	if args != nil && !args.ImpliedPort {
		argsPort := args.Port
		if argsPort != nil && *argsPort != 0 {
			port = *argsPort
		}
	}
	return port
}
