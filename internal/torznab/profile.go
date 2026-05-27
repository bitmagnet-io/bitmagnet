package torznab

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Profile struct {
	ID                      string `validate:"required"`
	Title                   string
	DisableOrderByRelevance bool
	DefaultLimit            uint
	MaxLimit                uint
	Tags                    []string
	BaseURL                 model.NullString
}

var ProfileDefault = Profile{
	ID:           "default",
	Title:        "bitmagnet",
	DefaultLimit: 100,
	MaxLimit:     100,
	BaseURL:      model.NullString{},
}

func (p Profile) MergeDefaults() Profile {
	if p.Title == "" {
		p.Title = ProfileDefault.Title
	}

	if p.DefaultLimit == 0 {
		p.DefaultLimit = ProfileDefault.DefaultLimit
	}

	if p.MaxLimit == 0 {
		p.MaxLimit = ProfileDefault.MaxLimit
	}

	if p.DefaultLimit > p.MaxLimit {
		p.DefaultLimit = p.MaxLimit
	}

	return p
}

func (p Profile) Caps() Caps {
	return Caps{
		Server: CapsServer{
			Title: p.Title,
		},
		Limits: CapsLimits{
			Max:     p.MaxLimit,
			Default: p.DefaultLimit,
		},
		Searching: CapsSearching{
			Search: CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					ParamQuery,
					ParamIMDBID,
					ParamTMDBID,
				}, ","),
			},
			TvSearch: CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					ParamQuery,
					ParamIMDBID,
					ParamTMDBID,
					ParamSeason,
					ParamEpisode,
				}, ","),
			},
			MovieSearch: CapsSearch{
				Available: "yes",
				SupportedParams: strings.Join([]string{
					ParamQuery,
					ParamIMDBID,
					ParamTMDBID,
				}, ","),
			},
			MusicSearch: CapsSearch{
				Available:       "yes",
				SupportedParams: ParamQuery,
			},
			AudioSearch: CapsSearch{
				Available: "no",
			},
			BookSearch: CapsSearch{
				Available:       "yes",
				SupportedParams: ParamQuery,
			},
		},
		Categories: CapsCategories{
			Categories: TopLevelCategories,
		},
	}
}
