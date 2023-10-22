package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/netip"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Config    Config
	Responder responder.Responder
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Server  Server
	AppHook fx.Hook `group:"app_hooks"`
}

type Server interface {
	Ready() <-chan struct{}
	Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (dht.RecvMsg, error)
}

func New(p Params) Result {
	s := &server{
		ready:            make(chan struct{}),
		stopped:          make(chan struct{}),
		localAddr:        netip.AddrPortFrom(netip.IPv4Unspecified(), 3334),
		socket:           NewSocket(),
		queries:          make(map[string]chan dht.RecvMsg),
		queryTimeout:     p.Config.QueryTimeout,
		queryLimiter:     concurrency.NewKeyedLimiter(rate.Every(time.Second*2), 2, 1000, time.Second*20),
		responder:        p.Responder,
		responderTimeout: time.Second * 5,
		idIssuer:         &variantIdIssuer{},
		logger:           p.Logger.Named("dht_server"),
	}
	return Result{
		Server: s,
		AppHook: fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := s.socket.Open(s.localAddr); err != nil {
					return fmt.Errorf("could not open socket: %w", err)
				}
				go s.start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				close(s.stopped)
				return nil
			},
		},
	}
}

type server struct {
	mutex            sync.Mutex
	ready            chan struct{}
	stopped          chan struct{}
	localAddr        netip.AddrPort
	socket           Socket
	queryLimiter     concurrency.KeyedLimiter
	queryTimeout     time.Duration
	queries          map[string]chan dht.RecvMsg
	responder        responder.Responder
	responderTimeout time.Duration
	idIssuer         IdIssuer
	logger           *zap.SugaredLogger
}

func (s *server) Ready() <-chan struct{} {
	return s.ready
}

func (s *server) start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.read(ctx)
	close(s.ready)
	<-s.stopped
}

func (s *server) read(ctx context.Context) {
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
		if ctx.Err() != nil {
			return
		}

		n, from, err := s.socket.Receive(buffer)
		if err != nil {
			// Socket is probably closed; if we're not shutting down then panic
			if ctx.Err() == nil {
				panic(fmt.Errorf("socket read error: %w", err))
			}
			return
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
			s.logger.Debugw("could not unmarshal packet data", "err", err)
			continue
		}

		recvMsg := dht.RecvMsg{
			Msg:  msg,
			From: from,
		}

		switch msg.Y {
		case dht.YQuery:
			go s.handleQuery(recvMsg)
		case dht.YResponse, dht.YError:
			go s.handleResponse(recvMsg)
		}
	}
}

func (s *server) handleQuery(msg dht.RecvMsg) {
	ctx, cancel := context.WithTimeout(context.Background(), s.responderTimeout)
	defer cancel()
	res := dht.Msg{
		T: msg.Msg.T,
		Y: dht.YResponse,
	}
	ret, retErr := s.responder.Respond(ctx, msg)
	if retErr != nil {
		dhtErr := &dht.Error{}
		if ok := errors.As(retErr, dhtErr); ok {
			res.E = dhtErr
		} else {
			res.E = &dht.Error{
				Code: dht.ErrorCodeServerError,
				Msg:  "server error",
			}
			s.logger.Errorw("server error", "msg", msg, "retErr", retErr)
		}
	} else {
		res.R = &ret
	}
	if sendErr := s.send(msg.From, res); sendErr != nil {
		s.logger.Debugw("could not send response", "msg", msg, "retErr", sendErr)
	}
}

func (s *server) handleResponse(msg dht.RecvMsg) {
	transactionId := msg.Msg.T
	s.mutex.Lock()
	ch, ok := s.queries[transactionId]
	s.mutex.Unlock()
	if ok {
		ch <- msg
	}
}

func (s *server) Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (r dht.RecvMsg, err error) {
	if limitErr := s.queryLimiter.Wait(ctx, addr.Addr().String()); limitErr != nil {
		return r, limitErr
	}
	transactionId := s.idIssuer.Issue()
	ch := make(chan dht.RecvMsg, 1)
	s.mutex.Lock()
	s.queries[transactionId] = ch
	s.mutex.Unlock()
	defer (func() {
		s.mutex.Lock()
		delete(s.queries, transactionId)
		s.mutex.Unlock()
	})()
	msg := dht.Msg{
		Q: q,
		T: transactionId,
		A: &args,
		Y: dht.YQuery,
	}
	if sendErr := s.send(addr, msg); sendErr != nil {
		s.logger.Debugw("could not send query", "msg", msg, "sendErr", sendErr)
		err = sendErr
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
			s.logger.Debugw("error response", "msg", msg, "res", res, "err", err)
		} else if r.Msg.R == nil {
			err = errors.New("return data missing from response")
		} else {
			s.logger.Debugw("successful response", "msg", msg, "res", res)
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
