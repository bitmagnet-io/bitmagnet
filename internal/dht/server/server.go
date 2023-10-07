package server

import (
	"context"
	"errors"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/anacrolix/dht/v2/transactions"
	"github.com/anacrolix/torrent/bencode"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/responder"
	sockaddr "github.com/libp2p/go-sockaddr/net"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
	"golang.org/x/time/rate"
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
	Start() error
	Shutdown() error
	Query(ctx context.Context, addr krpc.NodeAddr, q string, args krpc.MsgArgs) (dht.RecvMsg, error)
	// QueryUrgent is a query that is not rate limited (exposed for use by health check)
	QueryUrgent(ctx context.Context, addr krpc.NodeAddr, q string, args krpc.MsgArgs) (dht.RecvMsg, error)
}

func New(p Params) (r Result, err error) {
	localIp := [4]byte{0, 0, 0, 0}
	r.Server = &server{
		localIp: localIp,
		localAddr: krpc.NodeAddr{
			IP:   localIp[:],
			Port: 0,
		},
		queryTimeout:     p.Config.QueryTimeout,
		mutex:            &sync.Mutex{},
		limiter:          rate.NewLimiter(rate.Every(p.Config.RateLimit), 1),
		responder:        p.Responder,
		responderTimeout: time.Second * 5,
		logger:           p.Logger.Named("dht_server"),
	}
	r.AppHook = fx.Hook{
		OnStart: func(ctx context.Context) error {
			go (func() {
				// todo Handle error
				_ = r.Server.Start()
			})()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return r.Server.Shutdown()
		},
	}
	return
}

type server struct {
	mutex            *sync.Mutex
	started          bool
	stop             func()
	localIp          [4]byte
	localAddr        krpc.NodeAddr
	fd               int
	limiter          *rate.Limiter
	queryTimeout     time.Duration
	queries          map[string]chan dht.RecvMsg
	responder        responder.Responder
	responderTimeout time.Duration
	logger           *zap.SugaredLogger
}

var (
	ErrServerAlreadyStarted = errors.New("dht: Server already started")
	ErrServerClosed         = errors.New("dht: Server closed")
)

func (s *server) Start() error {
	s.mutex.Lock()
	unlockOnce := sync.Once{}
	unlock := func() {
		unlockOnce.Do(s.mutex.Unlock)
	}
	defer unlock()
	if s.started {
		return ErrServerAlreadyStarted
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fd, socketErr := unix.Socket(unix.SOCK_DGRAM, unix.AF_INET, 0)
	if socketErr != nil {
		return socketErr
	}
	if bindErr := unix.Bind(fd, &unix.SockaddrInet4{Addr: s.localIp, Port: s.localAddr.Port}); bindErr != nil {
		return bindErr
	}
	s.fd = fd
	s.queries = make(map[string]chan dht.RecvMsg)
	s.stop = cancel
	s.started = true
	unlock()
	go s.read(ctx)
	<-ctx.Done()
	if shutdownErr := s.Shutdown(); shutdownErr != nil {
		return shutdownErr
	}
	return ErrServerClosed
}

func (s *server) Shutdown() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.started {
		s.started = false
		s.stop()
		for k := range s.queries {
			delete(s.queries, k)
		}
		return unix.Close(s.fd)
	}
	return nil
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

		n, fromS, err := unix.Recvfrom(s.fd, buffer, 0)
		if err != nil {
			// Socket is probably closed, log an error if we're not shutting down
			if ctx.Err() == nil {
				s.logger.Errorw("socket read error", "err", err)
			}
			s.stop()
			return
		}

		if n == 0 {
			/* Datagram sockets in various domains  (e.g., the UNIX and Internet domains) permit
			 * zero-length datagrams. When such a datagram is received, the return value (n) is 0.
			 */
			continue
		}

		var msg krpc.Msg
		err = bencode.Unmarshal(buffer[:n], &msg)
		if err != nil {
			s.logger.Debugw("could not unmarshal packet data", "err", err)
			continue
		}

		from := sockaddr.SockaddrToUDPAddr(fromS)
		if from == nil {
			s.logger.Warn("dht transport SockaddrToUDPAddr: null")
			// do something else here?
			continue
		}

		recvMsg := dht.RecvMsg{
			Msg: msg,
			From: krpc.NodeAddr{
				IP:   from.IP,
				Port: from.Port,
			},
		}

		switch msg.Y {
		case krpc.YQuery:
			go s.handleQuery(recvMsg)
		case krpc.YResponse, krpc.YError:
			go s.handleResponse(recvMsg)
		}
	}
}

func (s *server) handleQuery(msg dht.RecvMsg) {
	ctx, cancel := context.WithTimeout(context.Background(), s.responderTimeout)
	defer cancel()
	res := krpc.Msg{
		T: msg.Msg.T,
		Y: krpc.YResponse,
	}
	ret, err := s.responder.Respond(ctx, msg)
	if err != nil {
		krpcErr := &krpc.Error{}
		if ok := errors.As(err, krpcErr); ok {
			res.E = krpcErr
		} else {
			res.E = &krpc.Error{
				Code: krpc.ErrorCodeServerError,
				Msg:  "server error",
			}
			s.logger.Errorw("server error", "msg", msg, "err", err)
		}
	} else {
		res.R = &ret
	}
	s.logger.Debugw("responding", "msg", msg, "ret", ret)
	_ = s.send(msg.From, res)
}

func (s *server) handleResponse(msg dht.RecvMsg) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	transactionId := msg.Msg.T
	ch, ok := s.queries[transactionId]
	if ok {
		ch <- msg
	}
}

func (s *server) Query(ctx context.Context, addr krpc.NodeAddr, q string, args krpc.MsgArgs) (r dht.RecvMsg, err error) {
	if limitErr := s.limiter.Wait(ctx); limitErr != nil {
		err = limitErr
		return
	}
	return s.QueryUrgent(ctx, addr, q, args)
}

func (s *server) QueryUrgent(ctx context.Context, addr krpc.NodeAddr, q string, args krpc.MsgArgs) (r dht.RecvMsg, err error) {
	transactionId := nextTransactionId()
	ch := make(chan dht.RecvMsg, 1)
	s.mutex.Lock()
	started := s.started
	if !started {
		s.mutex.Unlock()
		err = ErrServerClosed
		return
	}
	s.queries[transactionId] = ch
	s.mutex.Unlock()
	defer (func() {
		s.mutex.Lock()
		close(ch)
		delete(s.queries, transactionId)
		s.mutex.Unlock()
	})()
	msg := krpc.Msg{
		Q: q,
		T: transactionId,
		A: &args,
		Y: krpc.YQuery,
	}
	if sendErr := s.send(addr, msg); sendErr != nil {
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
		if res.Msg.Y == krpc.YError {
			err = res.Msg.E
			if err == nil {
				err = errors.New("error missing from response")
			}
			s.logger.Debugw("error response", "msg", msg, "res", res, "err", err)
		}
		return
	}
}

func (s *server) send(addr krpc.NodeAddr, msg krpc.Msg) (err error) {
	data, encodeErr := bencode.Marshal(msg)
	if encodeErr != nil {
		err = encodeErr
		return
	}
	addrS := sockaddr.NetAddrToSockaddr(addr.UDP())
	if addrS == nil {
		s.logger.Debugw("wrong net address for the remote peer", "addr", addr)
		err = errors.New("wrong net address for the remote peer")
		return
	}
	if sendErr := unix.Sendto(s.fd, data, 0, addrS); sendErr != nil {
		s.logger.Debugw("could not send packet to remote peer", "addr", addr, "err", err)
		err = sendErr
		return
	}
	return nil
}

func nextTransactionId() string {
	return transactions.DefaultIdIssuer.Issue()
}
