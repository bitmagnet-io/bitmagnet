package model

import (
	"fmt"
	"time"
)

var timeZero = time.Time{}.UnixMilli()

func (t *Torrent) EnsureTimes(now int64) {
	if t.GetCreatedAt() == timeZero {
		t.CreatedAt = now
	}

	if t.GetUpdatedAt() == timeZero {
		t.UpdatedAt = now
	}

	for _, tf := range t.GetFiles() {
		tf.EnsureTimes(now)
	}

	for _, ts := range t.GetSources() {
		ts.EnsureTimes(now)
	}
}

func (tf *TorrentFile) EnsureTimes(now int64) {
	if tf.GetCreatedAt() == timeZero {
		tf.CreatedAt = now
	}

	if tf.GetUpdatedAt() == timeZero {
		tf.UpdatedAt = now
	}
}

func (ts *TorrentSource) EnsureTimes(now int64) {
	if ts.GetCreatedAt() == timeZero {
		ts.CreatedAt = now
	}

	if ts.GetUpdatedAt() == timeZero {
		ts.UpdatedAt = now
	}
}

func (tc *TorrentContent) EnsureTimes(now int64) {
	if tc.GetCreatedAt() == timeZero {
		tc.CreatedAt = now
	}

	if tc.GetUpdatedAt() == timeZero {
		tc.UpdatedAt = now
	}

	if t := tc.GetTorrent(); t != nil {
		t.EnsureTimes(now)
	}

	if c := tc.GetContent(); tc != nil {
		c.EnsureTimes(now)
	}
}

func (c *Content) EnsureTimes(now int64) {
	if c.GetCreatedAt() == timeZero {
		c.CreatedAt = now
	}

	if c.GetUpdatedAt() == timeZero {
		c.UpdatedAt = now
	}

	for _, ca := range c.GetAttributes() {
		ca.EnsureTimes(now)
	}

	for _, cc := range c.GetCollections() {
		cc.EnsureTimes(now)
	}
}

func (r *ContentRef) String() string {
	return fmt.Sprintf("%s:%s:%s", r.GetType(), r.GetSource(), r.GetId())
}

func (cc *ContentCollection) EnsureTimes(now int64) {
	if cc.GetCreatedAt() == timeZero {
		cc.CreatedAt = now
	}

	if cc.GetUpdatedAt() == timeZero {
		cc.UpdatedAt = now
	}
}

func (ca *ContentAttribute) EnsureTimes(now int64) {
	if ca.GetCreatedAt() == timeZero {
		ca.CreatedAt = now
	}

	if ca.GetUpdatedAt() == timeZero {
		ca.UpdatedAt = now
	}
}
