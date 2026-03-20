package dhtcrawler

import (
	"context"
	"net"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func ResolveBootstrapAddr(strAddr string) (netip.AddrPort, error) {
	addr, err := net.ResolveUDPAddr("udp", strAddr)
	if err != nil {
		return netip.AddrPort{}, err
	}

	// Normalize IPv4-mapped IPv6 addresses so downstream DHT networking uses the IPv4 family.
	if ip4 := addr.IP.To4(); ip4 != nil {
		addr.IP = ip4
	}

	return addr.AddrPort(), nil
}

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, strAddr := range c.bootstrapNodes {
				addr, err := ResolveBootstrapAddr(strAddr)
				if err != nil {
					c.logger.Warnf("failed to resolve bootstrap node address: %s", err)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, addr):
					continue
				}
			}
		}

		interval = c.reseedBootstrapNodesInterval
	}
}
