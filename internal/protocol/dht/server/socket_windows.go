package server

import (
	"errors"
	"fmt"
	"net/netip"
	"syscall"

	"golang.org/x/sys/windows"
)

type socket struct {
	fd windows.Handle
}

func newSocket() Socket {
	fd, _ := windows.Socket(windows.AF_INET, windows.SOCK_DGRAM, 0)
	return &socket{fd: fd}
}

func (s *socket) Open(localAddr netip.AddrPort) error {
	sAddr, addrErr := addrPortToSockaddr(localAddr)
	if addrErr != nil {
		return addrErr
	}
	return windows.Bind(s.fd, sAddr)
}

func (s *socket) Close() error {
	return windows.Close(s.fd)
}

func (s *socket) Send(remoteAddr netip.AddrPort, data []byte) error {
	sAddr, addrErr := addrPortToSockaddr(remoteAddr)
	if addrErr != nil {
		return addrErr
	}
	return windows.Sendto(s.fd, data, 0, sAddr)
}

func (s *socket) Receive(data []byte) (int, netip.AddrPort, error) {
	n, sAddr, recvErr := windows.Recvfrom(s.fd, data, 0)
	if recvErr != nil {
		if ignoreReadFromError(recvErr) {
			return 0, netip.AddrPort{}, nil
		}
		return n, netip.AddrPort{}, recvErr
	}
	addr, addrErr := sockaddrToAddrPort(sAddr)
	if addrErr != nil {
		return n, netip.AddrPort{}, addrErr
	}
	return n, addr, nil
}

// see https://github.com/bitmagnet-io/bitmagnet/pull/203 and https://github.com/anacrolix/dht/issues/16
func ignoreReadFromError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case
			windows.WSAENETRESET,
			windows.WSAECONNRESET,
			windows.WSAECONNABORTED,
			windows.WSAECONNREFUSED,
			windows.WSAENETUNREACH,
			windows.WSAETIMEDOUT:
			return true
		}
	}
	return false
}

func addrPortToSockaddr(addr netip.AddrPort) (windows.Sockaddr, error) {
	if addr.Addr().Is4() {
		return &windows.SockaddrInet4{
			Addr: addr.Addr().As4(),
			Port: int(addr.Port()),
		}, nil
	}
	if addr.Addr().Is6() {
		return &windows.SockaddrInet6{
			Addr: addr.Addr().As16(),
			Port: int(addr.Port()),
		}, nil
	}
	return nil, errors.New("invalid address")
}

func sockaddrToAddrPort(addr windows.Sockaddr) (netip.AddrPort, error) {
	switch addr := addr.(type) {
	case *windows.SockaddrInet4:
		return netip.AddrPortFrom(netip.AddrFrom4(addr.Addr), uint16(addr.Port)), nil
	case *windows.SockaddrInet6:
		return netip.AddrPortFrom(netip.AddrFrom16(addr.Addr), uint16(addr.Port)), nil
	default:
		return netip.AddrPort{}, fmt.Errorf("unsupported sockaddr type: %T", addr)
	}
}
