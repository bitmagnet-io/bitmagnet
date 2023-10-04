package crawler

import (
	"context"
	"crypto/rand"
	adht "github.com/anacrolix/dht/v2"
	"github.com/anacrolix/dht/v2/krpc"
	"net"
	"time"
)

func (c *crawler) crawlBootstrapHosts(ctx context.Context) {
	for {
		for _, n := range adht.DefaultGlobalBootstrapHostPorts {
			go (func(node string) {
				c.logger.Debugf("bootstrap node: %s", node)
				t := [20]byte{}
				if _, randErr := rand.Read(t[:]); randErr != nil {
					c.logger.Errorw("could not generate random bytes", "err", randErr)
					return
				}
				addr, err := net.ResolveUDPAddr("udp", node)
				if err != nil {
					c.logger.Errorw("could not resolve UDP address of the bootstrapping node", "node", node)
					return
				}
				res, err := c.dhtServer.Query(ctx, krpc.NodeAddr{
					IP:   addr.IP.To4(),
					Port: addr.Port,
				}, "find_node", krpc.MsgArgs{
					Target: t,
				})
				if err == nil {
					peersToReceive := make([]krpc.NodeAddr, 0, len(res.Reply.R.Nodes))
					for _, thisNode := range res.Reply.R.Nodes {
						peersToReceive = append(peersToReceive, thisNode.Addr)
					}
					nAdded, nDiscarded := c.routingTable.ReceivePeers(peersToReceive...)
					c.logger.Debugw("received peers", "node", node, "nAdded", nAdded, "nDiscarded", nDiscarded)
				} else {
					c.logger.Debugw("find_node error", "node", node, "err", err)
				}
			})(n)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.crawlBootstrapHostsInterval):
			continue
		}
	}
}
