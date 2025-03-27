package search

import (
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func TorrentContentEpisodesCriteria(episodes model.Episodes) query.Criteria {
	return query.GenCriteria(func(query.DBContext) (query.Criteria, error) {
		and := make([]query.Criteria, 0, len(episodes))

		for _, s := range episodes.SeasonEntries() {
			if len(s.Episodes) == 0 {
				and = append(and, query.DBCriteria{
					SQL: fmt.Sprintf("torrent_contents.episodes #> '{%d}' = '{}'::jsonb", s.Season),
				})
			} else {
				keyParts := make([]string, 0, len(s.Episodes))
				for _, e := range s.Episodes {
					keyParts = append(keyParts, fmt.Sprintf("\"%d\":{}", e))
				}

				and = append(and, query.DBCriteria{
					SQL: fmt.Sprintf(
						"torrent_contents.episodes #> '{%d}' @> '{%s}'::jsonb",
						s.Season,
						strings.Join(keyParts, ","),
					),
				})
			}
		}

		return query.AndCriteria{
			Criteria: and,
		}, nil
	})
}
