package dhtcrawler

import (
	"context"
	"time"
)

func (c *crawler) rotateSoughtPeerId(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			c.soughtPeerId.Set(c.kTable.GeneratePeerID())
		}
	}
}
