package workflow

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
)

func (r workflow) Run(ctx context.Context, t model.Torrent) (classifier.Classification, error) {
	cl := classifier.Classification{}
	if !t.Hint.IsNil() {
		cl.ApplyHint(t.Hint)
	}
	// if possible, attach the existing content to the result to save some work:
	if !t.Hint.IsNil() && t.Hint.ContentSource.Valid {
		for _, tc := range t.Contents {
			if tc.ContentType.Valid &&
				tc.ContentType.ContentType == t.Hint.ContentType &&
				tc.ContentSource.Valid &&
				tc.ContentSource.String == t.Hint.ContentSource.String &&
				tc.ContentID.String == t.Hint.ContentID.String &&
				tc.Content.Source == tc.ContentSource.String {
				c := tc.Content
				cl.Content = &c
				break
			}
		}
	}
	exCtx := executionContext{
		Context:   ctx,
		torrent:   t,
		torrentPb: protobuf.NewTorrent(t),
		result:    cl,
	}
	return r.action.run(exCtx)
}
