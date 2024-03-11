//go:build !windows

package server

import (
	"errors"
	"fmt"
	"net/netip"

	"golang.org/x/sys/unix"
)

func newSocket() Socket {
	fd, sockErr := unix.Socket(unix.SOCK_DGRAM, unix.AF_INET, 0)
	if sockErr != nil {
		panic(fmt.Errorf("error creating socket: %w", sockErr))
	}
	return &socket{
		fd: fd,
	}
}

type socket struct {
	fd int
}

func (s *socket) Open(localAddr netip.AddrPort) error {
	sAddr, addrErr := addrPortToSockaddr(localAddr)
	if addrErr != nil {
		return addrErr
	}
	if bindErr := unix.Bind(s.fd, sAddr); bindErr != nil {
		return bindErr
	}
	return nil
}

func (s *socket) Close() error {
	return unix.Close(s.fd)
}

func (s *socket) Send(remoteAddr netip.AddrPort, data []byte) error {
	sAddr, addrErr := addrPortToSockaddr(remoteAddr)
	if addrErr != nil {
		return addrErr
	}
	return unix.Sendto(s.fd, data, 0, sAddr)
}

func (s *socket) Receive(data []byte) (int, netip.AddrPort, error) {
	n, sAddr, recvErr := unix.Recvfrom(s.fd, data, 0)
	if recvErr != nil {
		return n, netip.AddrPort{}, recvErr
	}
	addr, addrErr := sockaddrToAddrPort(sAddr)
	if addrErr != nil {
		return n, netip.AddrPort{}, addrErr
	}
	return n, addr, nil
}

func addrPortToSockaddr(addr netip.AddrPort) (unix.Sockaddr, error) {
	if addr.Addr().Is4() {
		return &unix.SockaddrInet4{
			Addr: addr.Addr().As4(),
			Port: int(addr.Port()),
		}, nil
	}
	if addr.Addr().Is6() {
		return &unix.SockaddrInet6{
			Addr: addr.Addr().As16(),
			Port: int(addr.Port()),
		}, nil
	}
	return nil, errors.New("invalid address")
}

func sockaddrToAddrPort(addr unix.Sockaddr) (netip.AddrPort, error) {
	switch addr := addr.(type) {
	case *unix.SockaddrInet4:
		return netip.AddrPortFrom(netip.AddrFrom4(addr.Addr), uint16(addr.Port)), nil
	case *unix.SockaddrInet6:
		return netip.AddrPortFrom(netip.AddrFrom16(addr.Addr), uint16(addr.Port)), nil
	default:
		return netip.AddrPort{}, fmt.Errorf("unsupported sockaddr type: %T", addr)
	}
}
