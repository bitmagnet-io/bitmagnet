package socket

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"sync"
	"syscall"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"golang.org/x/sys/windows"
)

func newSocketWindows() Runner {
	return &socketWindows{}
}

func init() {
	addAdapter("windows", newSocketWindows)
}

type socketWindows struct {
	localAddr netip.AddrPort
	mtx       sync.RWMutex
	fd        windows.Handle
}

func (s *socketWindows) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (shutdowner runner.Shutdowner, err error) {
		defer func() {
			if err != nil {
				cancel(err)
			}
		}()

		shutdowner = runner.NopShutdowner

		s.mtx.Lock()
		defer s.mtx.Unlock()

		if s.fd != 0 {
			return runner.NopShutdowner, fmt.Errorf("%w: %w: %w", Err, runner.ErrAlreadyRunning)
		}

		fd, err := windows.Socket(windows.AF_INET, windows.SOCK_DGRAM, windows.IPPROTO_UDP)
		if err != nil {
			return runner.NopShutdowner,
				fmt.Errorf("%w: %w: %w", Err, ErrCreateFailed, errors.Join(err, windows.Close(fd)))
		}

		s.fd = fd

		shutdowner = func(context.Context) error {
			s.mtx.Lock()
			defer s.mtx.Unlock()

			s.fd = 0

			err := windows.Close(fd)
			if err != nil {
				return fmt.Errorf("%w: %w: %w", Err, ErrCloseFailed, err)
			}

			return nil
		}

		return
	}
}

func (s *socketWindows) Send(remoteAddr netip.AddrPort, data []byte) error {
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

	err = windows.Sendto(fd, data, 0, sAddr)

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

func (s *socketWindows) Receive(data []byte) (int, netip.AddrPort, error) {
	s.mtx.RLock()
	fd := s.fd
	s.mtx.RUnlock()

	if fd == 0 {
		return 0, netip.AddrPort{}, fmt.Errorf("%w: %w", ErrReceiveFailed, ErrClosed)
	}

	n, sAddr, err := windows.Recvfrom(fd, data, 0)
	if err != nil {
		s.mtx.RLock()
		isShutdown := s.fd != fd
		s.mtx.RUnlock()

		if isShutdown {
			err = runner.ErrShutdownRequested
		}

		if ignoreReadFromError(err) {
			return 0, netip.AddrPort{}, nil
		}

		return n, netip.AddrPort{}, fmt.Errorf("%w: %w: %w", Err, ErrReceiveFailed, err)
	}

	addr, err := sockaddrToAddrPort(sAddr)
	if err != nil {
		return n, netip.AddrPort{}, fmt.Errorf("%w: %w: %w", Err, ErrReceiveFailed, err)
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

	return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, addr.String())
}

func sockaddrToAddrPort(addr windows.Sockaddr) (netip.AddrPort, error) {
	switch addr := addr.(type) {
	case *windows.SockaddrInet4:
		return netip.AddrPortFrom(netip.AddrFrom4(addr.Addr), uint16(addr.Port)), nil
	case *windows.SockaddrInet6:
		return netip.AddrPortFrom(netip.AddrFrom16(addr.Addr), uint16(addr.Port)), nil
	default:
		return netip.AddrPort{}, fmt.Errorf("%w: %T", ErrUnsupportedAddressType, addr)
	}
}
