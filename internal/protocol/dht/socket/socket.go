package socket

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Socket interface {
	Send(netip.AddrPort, []byte) error
	Receive([]byte) (int, netip.AddrPort, error)
}

type Runner interface {
	Socket
	runner.Provider
}
