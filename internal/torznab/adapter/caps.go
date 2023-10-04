package adapter

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"strings"
)

func (a adapter) Caps(context.Context) (torznab.Caps, error) {
	return torznab.Caps{
		Server: torznab.CapsServer{
			Title: a.title,
		},
		Limits: torznab.CapsLimits{
			Max:     a.maxLimit,
			Default: a.defaultLimit,
		},
		Searching: torznab.CapsSearching{
			Search: torznab.CapsSearch{
				Available:       "yes",
				SupportedParams: torznab.ParamQuery,
			},
			TvSearch: torznab.CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					torznab.ParamQuery,
					torznab.ParamImdbId,
					torznab.ParamSeason,
					torznab.ParamEpisode,
				}, ","),
			},
			MovieSearch: torznab.CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					torznab.ParamQuery,
					torznab.ParamImdbId,
				}, ","),
			},
			MusicSearch: torznab.CapsSearch{
				Available:       "yes",
				SupportedParams: torznab.ParamQuery,
			},
			AudioSearch: torznab.CapsSearch{
				Available: "no",
			},
			BookSearch: torznab.CapsSearch{
				Available:       "yes",
				SupportedParams: torznab.ParamQuery,
			},
		},
		Categories: torznab.CapsCategories{
			Categories: torznab.TopLevelCategories,
		},
	}, nil
}
