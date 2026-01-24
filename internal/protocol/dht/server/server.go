package server

import (
	"context"
	"errors"
	"net/netip"
	"sync"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

type Server interface {
	Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (dht.RecvMsg, error)
}

type Runner interface {
	Server
	runner.Provider
}

type server struct {
	mutex            sync.Mutex
	socket           socket.Socket
	queryTimeout     time.Duration
	queries          map[string]chan dht.RecvMsg
	responder        responder.Responder
	responderTimeout time.Duration
	idIssuer         IDIssuer
	logger           *zap.Logger
}

func (s *server) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		shutdown := make(chan struct{})

		go func() {
			err := s.read(ctx)

			select {
			case <-shutdown:
				if errors.Is(err, context.Canceled) {
					err = nil
				}
			default:
			}

			cancel(err)
		}()

		return func(context.Context) error {
			close(shutdown)
			cancel(runner.ErrShutdownRequested)

			return nil
		}, nil
	}
}

func (s *server) read(ctx context.Context) error {
	/*   The field size sets a theoretical limit of 65,535 bytes (8 byte header + 65,527 bytes of
	 * data) for a UDP datagram. However the actual limit for the data length, which is imposed by
	 * the underlying IPv4 protocol, is 65,507 bytes (65,535 − 8 byte UDP header − 20 byte IP
	 * header).
	 *
	 *   In IPv6 jumbograms it is possible to have UDP packets of size greater than 65,535 bytes.
	 * RFC 2675 specifies that the length field is set to zero if the length of the UDP header plus
	 * UDP data is greater than 65,535.
	 *
	 * https://en.wikipedia.org/wiki/User_Datagram_Protocol
	 */
	buffer := make([]byte, 65507)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, from, err := s.socket.Receive(buffer)
		if err != nil {
			if errors.Is(err, runner.ErrShutdownRequested) {
				return nil
			}

			return err
		}

		if n == 0 {
			/* Datagram sockets in various domains  (e.g., the UNIX and Internet domains) permit
			 * zero-length datagrams. When such a datagram is received, the return value (n) is 0.
			 */
			continue
		}

		var msg dht.Msg

		err = bencode.Unmarshal(buffer[:n], &msg)
		if err != nil {
			s.logger.Debug("failed to unmarshal packet data", zap.Error(err))
			continue
		}

		recvMsg := dht.RecvMsg{
			Msg:  msg,
			From: from,
		}

		switch msg.Y {
		case dht.YQuery:
			go s.handleQuery(ctx, recvMsg)
		case dht.YResponse, dht.YError:
			go s.handleResponse(recvMsg)
		}
	}
}

func (s *server) handleQuery(ctx context.Context, msg dht.RecvMsg) {
	ctx, cancel := context.WithTimeout(ctx, s.responderTimeout)
	defer cancel()

	res := dht.Msg{
		T: msg.Msg.T,
		Y: dht.YResponse,
	}

	ret, err := s.responder.Respond(ctx, msg)
	if err != nil {
		dhtErr := &dht.Error{}
		if ok := errors.As(err, dhtErr); ok {
			res.E = dhtErr
		} else {
			res.E = &dht.Error{
				Code: dht.ErrorCodeServerError,
				Msg:  "server error",
			}

			s.logger.Error("server error", zap.Any("msg", msg), zap.Error(err))
		}
	} else {
		res.R = &ret
	}

	if err = s.send(msg.From, res); err != nil {
		s.logger.Debug("failed to send response", zap.Any("msg", msg), zap.Error(err))
	}
}

func (s *server) handleResponse(msg dht.RecvMsg) {
	transactionID := msg.Msg.T

	s.mutex.Lock()
	ch, ok := s.queries[transactionID]
	s.mutex.Unlock()

	if ok {
		ch <- msg
	}
}

func (s *server) Query(
	ctx context.Context,
	addr netip.AddrPort,
	q string,
	args dht.MsgArgs,
) (r dht.RecvMsg, err error) {
	transactionID := s.idIssuer.Issue()
	ch := make(chan dht.RecvMsg, 1)

	s.mutex.Lock()
	s.queries[transactionID] = ch
	s.mutex.Unlock()

	defer (func() {
		s.mutex.Lock()
		delete(s.queries, transactionID)
		s.mutex.Unlock()
	})()

	msg := dht.Msg{
		Q: q,
		T: transactionID,
		A: &args,
		Y: dht.YQuery,
	}
	if err = s.send(addr, msg); err != nil {
		s.logger.Debug("failed to send query", zap.Any("msg", msg), zap.Error(err))

		return
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	select {
	case <-queryCtx.Done():
		err = queryCtx.Err()
		return
	case res, ok := <-ch:
		if !ok {
			err = errors.New("channel closed")
			return
		}

		r = res

		if res.Msg.Y == dht.YError {
			err = res.Msg.E
			if err == nil {
				err = errors.New("error missing from response")
			}
		} else if r.Msg.R == nil {
			err = errors.New("return data missing from response")
		}

		return
	}
}

func (s *server) send(addr netip.AddrPort, msg dht.Msg) error {
	data, encodeErr := bencode.Marshal(msg)
	if encodeErr != nil {
		return encodeErr
	}

	sendErr := s.socket.Send(addr, data)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
