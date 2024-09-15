package search

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"strings"
)

func TorrentInfoHashCriteria(infoHashes ...protocol.ID) query.Criteria {
	return infoHashCriteria(model.TableNameTorrent, infoHashes...)
}

func infoHashCriteria(table string, infoHashes ...protocol.ID) query.Criteria {
	if len(infoHashes) == 0 {
		return query.DbCriteria{
			Sql: "FALSE",
		}
	}
	decodes := make([]string, len(infoHashes))
	for i, infoHash := range infoHashes {
		decodes[i] = fmt.Sprintf("DECODE('%s', 'hex')", infoHash.String())
	}
	return query.DbCriteria{
		Sql: fmt.Sprintf("%s.info_hash IN (%s)", table, strings.Join(decodes, ", ")),
	}
}
