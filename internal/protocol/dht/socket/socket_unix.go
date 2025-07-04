//go:build !windows

package socket

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"golang.org/x/sys/unix"
)

func newSocketUnix(localAddr netip.AddrPort) Runner {
	return &socketUnix{
		localAddr: localAddr,
	}
}

func init() {
	addAdapter("unix", newSocketUnix)
}

type socketUnix struct {
	localAddr netip.AddrPort
	mtx       sync.RWMutex
	fd        int
}

func (s *socketUnix) Runner(
	_ context.Context,
	cancel context.CancelCauseFunc,
) (shutdowner runner.Shutdowner, err error) {
	defer func() {
		if err != nil {
			cancel(err)
		}
	}()

	shutdowner = runner.NopShutdowner

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.fd != 0 {
		err = fmt.Errorf("%w: %w: %w", Err, ErrOpenFailed, runner.ErrAlreadyRunning)
		return
	}

	sAddr, err := addrPortToSockaddr(s.localAddr)
	if err != nil {
		err = fmt.Errorf("%w: %w: %w", Err, ErrOpenFailed, err)
		return
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		err = fmt.Errorf("%w: %w: %w: %w", Err, ErrOpenFailed, ErrCreateFailed, err)
		return
	}

	// Avoid address already in use error when restarting:
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		err = fmt.Errorf(
			"%w: %w: %w: %w",
			Err,
			ErrOpenFailed,
			ErrSetOptionFailed,
			errors.Join(err, unix.Close(fd)),
		)

		return
	}

	err = unix.Bind(fd, sAddr)
	if err != nil {
		err = fmt.Errorf(
			"%w: %w: %w: %w",
			Err,
			ErrOpenFailed,
			ErrBindFailed,
			errors.Join(err, unix.Close(fd)),
		)

		return
	}

	s.fd = fd

	shutdowner = func(context.Context) error {
		s.mtx.Lock()
		defer s.mtx.Unlock()

		s.fd = 0

		err := unix.Close(fd)
		if err != nil {
			return fmt.Errorf("%w: %w: %w", Err, ErrCloseFailed, err)
		}

		return nil
	}

	return
}

func (s *socketUnix) Send(remoteAddr netip.AddrPort, data []byte) error {
	sAddr, err := addrPortToSockaddr(remoteAddr)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrSendFailed, err)
	}

	s.mtx.RLock()
	fd := s.fd
	s.mtx.RUnlock()

	if fd == 0 {
		return fmt.Errorf("%w: %w", ErrSendFailed, ErrClosed)
	}

	err = unix.Sendto(fd, data, 0, sAddr)
	if err != nil {
		s.mtx.RLock()
		isShutdown := s.fd != fd
		s.mtx.RUnlock()

		if isShutdown {
			err = runner.ErrShutdownRequested
		}

		return fmt.Errorf("%w: %w: %w", Err, ErrSendFailed, err)
	}

	return nil
}

func (s *socketUnix) Receive(data []byte) (int, netip.AddrPort, error) {
	s.mtx.RLock()
	fd := s.fd
	s.mtx.RUnlock()

	if fd == 0 {
		return 0, netip.AddrPort{}, fmt.Errorf("%w: %w", ErrReceiveFailed, ErrClosed)
	}

	n, sAddr, err := unix.Recvfrom(fd, data, 0)
	if err != nil {
		s.mtx.RLock()
		isShutdown := s.fd != fd
		s.mtx.RUnlock()

		if isShutdown {
			err = runner.ErrShutdownRequested
		}

		return 0, netip.AddrPort{}, fmt.Errorf("%w: %w: %w", Err, ErrReceiveFailed, err)
	}

	addr, err := sockaddrToAddrPort(sAddr)
	if err != nil {
		return n, netip.AddrPort{}, fmt.Errorf("%w: %w: %w", Err, ErrReceiveFailed, err)
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

	return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, addr.String())
}

func sockaddrToAddrPort(addr unix.Sockaddr) (netip.AddrPort, error) {
	switch addr := addr.(type) {
	case *unix.SockaddrInet4:
		return netip.AddrPortFrom(netip.AddrFrom4(addr.Addr), uint16(addr.Port)), nil
	case *unix.SockaddrInet6:
		return netip.AddrPortFrom(netip.AddrFrom16(addr.Addr), uint16(addr.Port)), nil
	default:
		return netip.AddrPort{}, fmt.Errorf("%w: %T", ErrUnsupportedAddressType, addr)
	}
}
