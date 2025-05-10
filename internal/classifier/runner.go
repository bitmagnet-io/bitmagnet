package classifier

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/common/types/ref"
)

type runner struct {
	dependencies
	flagDefinitions
	compiledFlags
	workflows map[string]action
}

func (r runner) Run(ctx context.Context, workflow string, flags Flags, t model.Torrent) (classification.Result, error) {
	w, ok := r.workflows[workflow]
	if !ok {
		return classification.Result{}, fmt.Errorf("workflow not found: %s", workflow)
	}

	cfs := make(map[string]ref.Val, len(r.flagDefinitions))

	for k, d := range r.flagDefinitions {
		if runtimeRawVal, ok := flags[k]; ok {
			rcf, err := d.celVal(runtimeRawVal)
			if err != nil {
				return classification.Result{}, fmt.Errorf(
					"invalid value for runtime flag '%s': %w",
					k,
					err,
				)
			}

			cfs[k] = rcf
		} else {
			cfs[k] = r.compiledFlags[k]
		}
	}

	cl := classification.Result{}
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
				content := tc.Content
				cl.AttachContent(&content)

				break
			}
		}
	}

	exCtx := executionContext{
		Context:      ctx,
		dependencies: r.dependencies,
		workflows:    r.workflows,
		flags:        cfs,
		torrent:      t,
		torrentPb:    protobuf.NewTorrent(t),
		result:       cl,
	}

	return w.run(exCtx)
}
