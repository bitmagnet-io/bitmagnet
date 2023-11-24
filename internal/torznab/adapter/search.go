package adapter

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"strconv"
	"strings"
)

func (a adapter) Search(ctx context.Context, req torznab.SearchRequest) (torznab.SearchResult, error) {
	options := []query.Option{search.TorrentContentDefaultOption()}
	if reqOptions, reqErr := a.searchRequestOptions(req); reqErr != nil {
		return torznab.SearchResult{}, reqErr
	} else {
		options = append(options, reqOptions...)
	}
	searchResult, searchErr := a.search.TorrentContent(ctx, options...)
	if searchErr != nil {
		return torznab.SearchResult{}, searchErr
	}
	return a.transformSearchResult(req, searchResult), nil
}

func (a adapter) searchRequestOptions(r torznab.SearchRequest) ([]query.Option, error) {
	var options []query.Option
	switch r.Type {
	case torznab.FunctionSearch:
		break
	case torznab.FunctionMovie:
		options = append(options, query.Where(search.ContentTypeCriteria(model.ContentTypeMovie)))
	case torznab.FunctionTv:
		options = append(options, query.Where(search.ContentTypeCriteria(model.ContentTypeTvShow)))
	case torznab.FunctionMusic:
		options = append(options, query.Where(search.ContentTypeCriteria(model.ContentTypeMusic)))
	case torznab.FunctionBook:
		options = append(options, query.Where(search.ContentTypeCriteria(model.ContentTypeBook)))
	default:
		return nil, torznab.Error{
			Code:        202,
			Description: fmt.Sprintf("no such function (%s)", r.Type),
		}
	}
	if r.Query != "" {
		options = append(options, query.QueryString(r.Query), query.OrderByQueryStringRank())
	}
	var catsCriteria []query.Criteria
	for _, cat := range r.Cats {
		var catCriteria []query.Criteria
		if torznab.CategoryMovies.Has(cat) {
			if r.Type != torznab.FunctionMovie {
				catCriteria = append(catCriteria, search.ContentTypeCriteria(model.ContentTypeMovie))
			}
			if torznab.CategoryMoviesSD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(model.VideoResolutionV480p))
			} else if torznab.CategoryMoviesHD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(
					model.VideoResolutionV720p,
					model.VideoResolutionV1080p,
					model.VideoResolutionV1440p,
					model.VideoResolutionV2160p,
				))
			} else if torznab.CategoryMoviesUHD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(model.VideoResolutionV2160p))
			} else if torznab.CategoryMovies3D.ID == cat {
				catCriteria = append(catCriteria, search.Video3dCriteria(
					model.Video3dV3D,
					model.Video3dV3DSBS,
					model.Video3dV3DOU,
				))
			}
		} else if torznab.CategoryTV.Has(cat) {
			if r.Type != torznab.FunctionTv {
				catCriteria = append(catCriteria, search.ContentTypeCriteria(model.ContentTypeTvShow))
			}
			if torznab.CategoryTVSD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(model.VideoResolutionV480p))
			} else if torznab.CategoryTVHD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(
					model.VideoResolutionV720p,
					model.VideoResolutionV1080p,
					model.VideoResolutionV1440p,
					model.VideoResolutionV2160p,
				))
			} else if torznab.CategoryTVUHD.ID == cat {
				catCriteria = append(catCriteria, search.VideoResolutionCriteria(model.VideoResolutionV2160p))
			}
		} else if torznab.CategoryAudio.Has(cat) {
			if r.Type != torznab.FunctionMusic {
				catCriteria = append(catCriteria, search.ContentTypeCriteria(model.ContentTypeMusic))
			}
		} else if torznab.CategoryBooks.Has(cat) {
			if r.Type != torznab.FunctionBook {
				catCriteria = append(catCriteria, search.ContentTypeCriteria(model.ContentTypeBook))
			}
		}
		if len(catCriteria) > 0 {
			catsCriteria = append(catsCriteria, query.And(catCriteria...))
		}
	}
	if len(catsCriteria) > 0 {
		options = append(options, query.Where(query.Or(catsCriteria...)))
	}
	if r.ImdbId.Valid {
		imdbId := r.ImdbId.String
		if !strings.HasPrefix(imdbId, "tt") {
			imdbId = "tt" + imdbId
		}
		var refs []model.ContentRef
		if r.Type != torznab.FunctionTv {
			refs = append(refs, model.ContentRef{
				Type:   model.ContentTypeMovie,
				Source: "imdb",
				ID:     imdbId,
			})
		}
		if r.Type != torznab.FunctionMovie {
			refs = append(refs, model.ContentRef{
				Type:   model.ContentTypeTvShow,
				Source: "imdb",
				ID:     imdbId,
			})
		}
		options = append(options, query.Where(search.ContentIdentifierCriteria(refs...)))
	}
	limit := a.defaultLimit
	if r.Limit.Valid {
		limit = r.Limit.Uint
		if limit > a.maxLimit {
			limit = a.maxLimit
		}
	}
	options = append(options, query.Limit(limit))
	if r.Offset.Valid {
		options = append(options, query.Offset(r.Offset.Uint))
	}
	// todo: Season and episodes
	return options, nil
}

func (a adapter) transformSearchResult(req torznab.SearchRequest, res search.TorrentContentResult) torznab.SearchResult {
	entries := make([]torznab.SearchResultItem, 0, len(res.Items))
	for _, item := range res.Items {
		category := "Unknown"
		if item.ContentType.Valid {
			category = item.ContentType.ContentType.Label()
		}
		date := item.Torrent.Sources[0].PublishedAt
		if date.IsZero() {
			date = item.CreatedAt
		}
		categoryId := torznab.CategoryOther.ID
		if item.ContentType.Valid {
			switch item.ContentType.ContentType {
			case model.ContentTypeMovie:
				categoryId = torznab.CategoryMovies.ID
			case model.ContentTypeTvShow:
				categoryId = torznab.CategoryTV.ID
			case model.ContentTypeMusic:
				categoryId = torznab.CategoryAudio.ID
			case model.ContentTypeBook:
				categoryId = torznab.CategoryBooks.ID
			case model.ContentTypeSoftware:
				categoryId = torznab.CategoryPC.ID
			case model.ContentTypeGame:
				categoryId = torznab.CategoryPCGames.ID
			}
		}
		attrs := []torznab.SearchResultItemTorznabAttr{
			{
				AttrName:  torznab.AttrCategory,
				AttrValue: strconv.Itoa(categoryId),
			},
			{
				AttrName:  torznab.AttrSize,
				AttrValue: strconv.FormatUint(item.Torrent.Size, 10),
			},
			{
				AttrName:  torznab.AttrPublishDate,
				AttrValue: date.Format(torznab.RssDateDefaultFormat),
			},
		}
		if seeders := item.Torrent.Seeders(); seeders.Valid {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrSeeders,
				AttrValue: strconv.Itoa(int(seeders.Uint)),
			})
		}
		if leechers := item.Torrent.Leechers(); leechers.Valid {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrPeers,
				AttrValue: strconv.Itoa(int(leechers.Uint)),
			})
		}
		if len(item.Torrent.Files) > 0 {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrFiles,
				AttrValue: strconv.Itoa(len(item.Torrent.Files)),
			})
		}
		if !item.ReleaseYear.IsNil() {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrYear,
				AttrValue: strconv.Itoa(int(item.ReleaseYear)),
			})
		}
		if len(item.Episodes) > 0 {
			// should we be adding all seasons and episodes here?
			seasons := item.Episodes.SeasonEntries()
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrSeason,
				AttrValue: strconv.Itoa(seasons[0].Season),
			})
			if len(seasons[0].Episodes) > 0 {
				attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
					AttrName:  torznab.AttrEpisode,
					AttrValue: strconv.Itoa(seasons[0].Episodes[0]),
				})
			}
		}
		if item.VideoCodec.Valid {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrVideo,
				AttrValue: item.VideoCodec.VideoCodec.Label(),
			})
		}
		if item.VideoResolution.Valid {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrResolution,
				AttrValue: item.VideoResolution.VideoResolution.Label(),
			})
		}
		if item.ReleaseGroup.Valid {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrTeam,
				AttrValue: item.ReleaseGroup.String,
			})
		}
		if imdbId, ok := item.Content.Identifier("imdb"); ok {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrImdb,
				AttrValue: imdbId[2:],
			})
		}
		entries = append(entries, torznab.SearchResultItem{
			Title:    item.Torrent.Name,
			Size:     item.Torrent.Size,
			Category: category,
			GUID:     item.InfoHash.String(),
			PubDate:  torznab.RssDate(date),
			Enclosure: torznab.SearchResultItemEnclosure{
				URL:    item.Torrent.MagnetUri(),
				Type:   "application/x-bittorrent;x-scheme-handler/magnet",
				Length: strconv.FormatUint(item.Torrent.Size, 10),
			},
			TorznabAttrs: attrs,
		})
	}
	return torznab.SearchResult{
		Channel: torznab.SearchResultChannel{
			Title: a.title,
			Response: torznab.SearchResultResponse{
				Offset: req.Offset.Uint,
				Total:  res.TotalCount,
			},
			Items: entries,
		},
	}
}
