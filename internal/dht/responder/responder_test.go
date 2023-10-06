package responder

import (
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	peerId := dht.RandomPeerID()
	r := responder{peerID: peerId}
	msg := krpc.Msg{
		Q: "ping",
	}
	ret, err := r.Respond(nil, msg)
	assert.NoError(t, err)
	assert.Equal(t, krpc.Return{ID: peerId}, ret)
}
