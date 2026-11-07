package classifier

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/common/types"
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

	// A configured MatchLab is authoritative for torrent-to-IMDb identification.
	// Runtime queue flags may disable API calls globally, but must not bypass
	// the provider and silently re-enable local/TMDB title searches.
	if r.matchlabClient != nil {
		cfs["matchlab_enabled"] = types.Bool(true)
		cfs["local_search_enabled"] = types.Bool(false)
	} else {
		// Queued jobs can survive a config change/restart. Never trust a
		// stale runtime matchlab_enabled=true flag after removing MATCHLAB_URL:
		// it would disable both MatchLab and the native search path.
		cfs["matchlab_enabled"] = types.Bool(false)
	}

	cl := classification.Result{}
	if !t.Hint.IsNil() {
		cl.ApplyHint(t.Hint)
	}
	// Previously attached movie/TV identities can be stale. MatchLab must
	// re-identify these releases rather than trusting cached hints.
	matchlabOwnsHint := r.matchlabClient != nil &&
		(t.Hint.ContentType == model.ContentTypeMovie || t.Hint.ContentType == model.ContentTypeTvShow)
	if !matchlabOwnsHint && !t.Hint.IsNil() && t.Hint.ContentSource.Valid {
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
