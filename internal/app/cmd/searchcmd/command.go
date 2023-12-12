// This command isn't currently intended to be usable, it's more of a testbed for trying things out, but may become user-friendly in future.

package searchcmd

import (
	"encoding/json"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Search search.Search
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	cmd := &cli.Command{
		Name: "search",
		Subcommands: []*cli.Command{
			{
				Name: "torrents",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "query",
					},
					&cli.UintFlag{
						Name:  "limit",
						Value: 10,
					},
					&cli.UintFlag{
						Name: "offset",
					},
					&cli.StringFlag{
						Name: "releaseDate",
					},
				},
				Action: func(ctx *cli.Context) error {
					result, searchErr := p.Search.TorrentContent(
						ctx.Context,
						search.TorrentContentDefaultOption(),
						query.QueryString(ctx.String("query")),
						//search.Where(
						//	search.ContentReleaseDateCriteriaString(ctx.String("releaseDate")),
						//),
						query.Limit(ctx.Uint("limit")),
						query.Offset(ctx.Uint("offset")),
						query.WithFacet(
							//search.ReleaseYearFacet(
							//	query.FacetHasFilter(query.FacetFilter{
							//		"2022": {},
							//		//"null": {},
							//	}),
							//	query.FacetIsAggregated(),
							//),
							//search.Video3dFacet(
							//	query.FacetIsAggregated(),
							//),
							//search.VideoCodecFacet(
							//	query.FacetIsAggregated(),
							//),
							//search.VideoModifierFacet(
							//	query.FacetIsAggregated(),
							//),
							//search.VideoResolutionFacet(
							//	query.FacetIsAggregated(),
							//),
							//search.VideoSourceFacet(
							//	query.FacetIsAggregated(),
							//),
							search.TorrentContentGenreFacet(
								query.FacetHasFilter(query.FacetFilter{
									"tmdb:10751": {},
									"tmdb:14":    {},
								}),
								//query.FacetIsAggregated(),
							),
						),
						query.OrderByQueryStringRank(),
						//query.Filter(query.FacetFilter{
						//	search.MovieGenreAggregatorKey: {
						//		//"tmdb:10751": {},
						//		//"tmdb:14": {},
						//		//"tmdb:35":    {},
						//	},
						//}),
						//search.OrderByColumn("torrents.created_at", true),
					)
					if searchErr != nil {
						return searchErr
					}
					jsonResult, jsonErr := json.Marshal(result)
					if jsonErr != nil {
						return jsonErr
					}
					p.Logger.Infof("Result: %v", string(jsonResult))
					return nil
				},
			},
		},
	}
	return Result{Command: cmd}, nil
}
