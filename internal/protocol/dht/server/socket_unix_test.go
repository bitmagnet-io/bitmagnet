//go:build !windows

package server

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

func TestAddrPortToSockaddr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr netip.AddrPort
		want any
	}{
		{
			name: "IPv4-mapped IPv6 is treated as IPv4",
			// ::ffff:1.2.3.4 — what Go's resolver can return on dual-stack
			// systems (https://github.com/bitmagnet-io/bitmagnet/issues/437).
			addr: netip.AddrPortFrom(netip.MustParseAddr("::ffff:1.2.3.4"), 6881),
			want: (*unix.SockaddrInet4)(nil),
		},
		{
			name: "pure IPv4",
			addr: netip.AddrPortFrom(netip.MustParseAddr("1.2.3.4"), 6881),
			want: (*unix.SockaddrInet4)(nil),
		},
		{
			name: "pure IPv6",
			addr: netip.AddrPortFrom(netip.MustParseAddr("2001:db8::1"), 6881),
			want: (*unix.SockaddrInet6)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := addrPortToSockaddr(tt.addr)
			assert.NoError(t, err)
			assert.IsType(t, tt.want, got)
		})
	}
}
