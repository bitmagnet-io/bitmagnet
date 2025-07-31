package socket_test

import (
	"context"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/bootstrap"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapters(t *testing.T) {
	t.Parallel()

	adapterNames := socket.AdapterNames()
	require.NotEmpty(t, adapterNames)

	idIssuer := &server.VariantIDIssuer{}

	for i, name := range adapterNames {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			t.Cleanup(cancel)

			adapter, err := socket.GetAdapter(name)
			require.NoError(t, err)

			sock := adapter(netip.AddrPortFrom(
				netip.IPv4Unspecified(),
				uint16(3334+i),
			))

			ctx, runnerCancel := context.WithCancelCause(ctx)
			t.Cleanup(func() {
				runnerCancel(nil)
			})

			shutdown, err := sock.Runner()(ctx, runnerCancel)
			require.NoError(t, err)

			t.Cleanup(func() {
				require.NoError(t, shutdown(t.Context()))
			})

			buffer := make([]byte, 65507)

			msgs := make(chan dht.Msg)

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					n, _, err := sock.Receive(buffer)
					if err != nil {
						continue
					}

					var msg dht.Msg

					err = bencode.Unmarshal(buffer[:n], &msg)
					if err != nil {
						continue
					}

					if msg.Y == dht.YResponse {
						select {
						case <-ctx.Done():
							return
						case msgs <- msg:
						}
					}
				}
			}()

			nodeID := protocol.RandomNodeID()

			expectedIDs := make(map[string]string)

			for _, nodeAddr := range bootstrap.DefaultBootstrapNodes {
				transactionID := idIssuer.Issue()

				data, err := bencode.Marshal(dht.Msg{
					Q: dht.QPing,
					T: transactionID,
					A: &dht.MsgArgs{ID: nodeID},
					Y: dht.YQuery,
				})
				require.NoError(t, err)

				addr, err := net.ResolveUDPAddr("udp", nodeAddr)
				if err != nil {
					t.Logf("could not resolve node address: %s", nodeAddr)
					continue
				}

				err = sock.Send(addr.AddrPort(), data)
				if err != nil {
					continue
				}

				expectedIDs[transactionID] = nodeAddr
			}

			receivedResponse := false

		outer:
			for {
				select {
				case <-ctx.Done():
					break outer
				case msg := <-msgs:
					nodeAddr, ok := expectedIDs[msg.T]
					if ok {
						t.Logf("received response from node: %s", nodeAddr)
						receivedResponse = true
					}
				}
			}

			assert.True(t, receivedResponse)
		})
	}
}
