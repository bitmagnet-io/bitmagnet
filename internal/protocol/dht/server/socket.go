package server

import (
	"net/netip"
)

type Socket interface {
	Open(localAddr netip.AddrPort) error
	Close() error
	Send(netip.AddrPort, []byte) error
	Receive([]byte) (int, netip.AddrPort, error)
}

func NewSocket(ip_type int) Socket {
	return newSocket(ip_type)
}
