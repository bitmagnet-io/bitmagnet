package socket

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func newSocketNet(localAddr netip.AddrPort) Runner {
	return &socketNet{
		localAddr: localAddr,
	}
}

func init() {
	addAdapter("net", newSocketNet)
}

type socketNet struct {
	localAddr netip.AddrPort
	mtx       sync.RWMutex
	conn      *net.UDPConn
}

func (s *socketNet) Runner() runner.Runner {
	return func(
		_ context.Context,
		cancel context.CancelCauseFunc,
	) (shutdowner runner.Shutdowner, err error) {
		defer func() {
			if err != nil {
				cancel(err)
			}
		}()

		s.mtx.Lock()
		defer s.mtx.Unlock()

		if s.conn != nil {
			err = fmt.Errorf("%w: %w: %w", Err, ErrOpenFailed, runner.ErrAlreadyRunning)
			return
		}

		conn, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(s.localAddr))
		if err != nil {
			err = fmt.Errorf("%w: %w: %w", Err, ErrOpenFailed, err)
			return
		}

		s.conn = conn

		shutdowner = func(context.Context) error {
			s.mtx.Lock()
			defer s.mtx.Unlock()

			err := s.conn.Close()
			s.conn = nil

			return err
		}

		return
	}
}

func (s *socketNet) Send(remoteAddr netip.AddrPort, data []byte) error {
	s.mtx.RLock()
	conn := s.conn
	s.mtx.RUnlock()

	if conn == nil {
		return ErrClosed
	}

	udpAddr := net.UDPAddrFromAddrPort(remoteAddr)

	_, err := conn.WriteToUDP(data, udpAddr)
	if err != nil {
		s.mtx.RLock()
		isShutdown := s.conn != conn
		s.mtx.RUnlock()

		if isShutdown {
			err = runner.ErrShutdownRequested
		}

		return fmt.Errorf("%w: %w: %w", Err, ErrSendFailed, err)
	}

	return err
}

func (s *socketNet) Receive(data []byte) (int, netip.AddrPort, error) {
	s.mtx.RLock()
	conn := s.conn
	s.mtx.RUnlock()

	if conn == nil {
		return 0, netip.AddrPort{}, ErrClosed
	}

	n, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		s.mtx.RLock()
		isShutdown := s.conn != conn
		s.mtx.RUnlock()

		if isShutdown {
			err = runner.ErrShutdownRequested
		}

		return 0, netip.AddrPort{}, fmt.Errorf("%w: %w: %w", Err, ErrReceiveFailed, err)
	}

	return n, addr.AddrPort(), nil
}
