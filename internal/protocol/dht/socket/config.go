package socket

import (
	"net/netip"
)

type Config struct {
	Adapter   string
	LocalAddr netip.AddrPort
}
