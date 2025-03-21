package adapter

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"strings"
)

func (a adapter) Caps(_ context.Context, profile torznab.Profile) (torznab.Caps, error) {
	return torznab.Caps{
		Server: torznab.CapsServer{
			Title: profile.Title,
		},
		Limits: torznab.CapsLimits{
			Max:     profile.MaxLimit,
			Default: profile.DefaultLimit,
		},
		Searching: torznab.CapsSearching{
			Search: torznab.CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					torznab.ParamQuery,
					torznab.ParamImdbId,
					torznab.ParamTmdbId,
				}, ","),
			},
			TvSearch: torznab.CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					torznab.ParamQuery,
					torznab.ParamImdbId,
					torznab.ParamTmdbId,
					torznab.ParamSeason,
					torznab.ParamEpisode,
				}, ","),
			},
			MovieSearch: torznab.CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					torznab.ParamQuery,
					torznab.ParamImdbId,
					torznab.ParamTmdbId,
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
