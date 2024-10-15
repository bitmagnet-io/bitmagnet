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
	options := []query.Option{search.TorrentContentDefaultOption(), query.WithTotalCount(false)}
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
		options = append(options, query.Where(search.TorrentContentTypeCriteria(model.ContentTypeMovie)))
	case torznab.FunctionTv:
		options = append(options, query.Where(search.TorrentContentTypeCriteria(model.ContentTypeTvShow)))
		if r.Season.Valid {
			episodes := make(model.Episodes)
			if r.Episode.Valid {
				episodes = episodes.AddEpisode(r.Season.Int, r.Episode.Int)
			} else {
				episodes = episodes.AddSeason(r.Season.Int)
			}
			options = append(options, query.Where(search.TorrentContentEpisodesCriteria(episodes)))
		}
	case torznab.FunctionMusic:
		options = append(options, query.Where(search.TorrentContentTypeCriteria(model.ContentTypeMusic)))
	case torznab.FunctionBook:
		options = append(options, query.Where(search.TorrentContentTypeCriteria(
			model.ContentTypeEbook,
			model.ContentTypeComic,
			model.ContentTypeAudiobook,
		)))
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
				catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeMovie))
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
				catCriteria = append(catCriteria, search.Video3DCriteria(
					model.Video3DV3D,
					model.Video3DV3DSBS,
					model.Video3DV3DOU,
				))
			}
		} else if torznab.CategoryTV.Has(cat) {
			if r.Type != torznab.FunctionTv {
				catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeTvShow))
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
		} else if torznab.CategoryXXX.Has(cat) {
			catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeXxx))
		} else if torznab.CategoryPC.Has(cat) {
			catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeSoftware, model.ContentTypeGame))
		} else if torznab.CategoryAudioAudiobook.Has(cat) {
			catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeAudiobook))
		} else if torznab.CategoryAudio.Has(cat) {
			catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeMusic))
		} else if torznab.CategoryBooksComics.Has(cat) {
			catCriteria = append(catCriteria, search.TorrentContentTypeCriteria(model.ContentTypeComic))
		} else if torznab.CategoryBooks.Has(cat) {
			options = append(options, query.Where(search.TorrentContentTypeCriteria(
				model.ContentTypeEbook,
				model.ContentTypeComic,
				model.ContentTypeAudiobook,
			)))
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
		var ct model.ContentType
		if r.Type != torznab.FunctionTv {
			ct = model.ContentTypeMovie
		} else if r.Type != torznab.FunctionMovie {
			ct = model.ContentTypeTvShow
		}
		options = append(options, query.Where(search.ContentAlternativeIdentifierCriteria(model.ContentRef{
			Type:   ct,
			Source: "imdb",
			ID:     imdbId,
		})))
	}
	if r.TmdbId.Valid {
		tmdbId := r.TmdbId.String
		var ct model.ContentType
		if r.Type != torznab.FunctionTv {
			ct = model.ContentTypeMovie
		} else if r.Type != torznab.FunctionMovie {
			ct = model.ContentTypeTvShow
		}
		options = append(options, query.Where(search.ContentCanonicalIdentifierCriteria(model.ContentRef{
			Type:   ct,
			Source: "tmdb",
			ID:     tmdbId,
		})))
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
	return options, nil
}

func (a adapter) transformSearchResult(req torznab.SearchRequest, res search.TorrentContentResult) torznab.SearchResult {
	entries := make([]torznab.SearchResultItem, 0, len(res.Items))
	for _, item := range res.Items {
		category := "Unknown"
		if item.ContentType.Valid {
			category = item.ContentType.ContentType.Label()
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
			case model.ContentTypeEbook:
				categoryId = torznab.CategoryBooks.ID
			case model.ContentTypeComic:
				categoryId = torznab.CategoryBooksComics.ID
			case model.ContentTypeAudiobook:
				categoryId = torznab.CategoryAudioAudiobook.ID
			case model.ContentTypeSoftware:
				categoryId = torznab.CategoryPC.ID
			case model.ContentTypeGame:
				categoryId = torznab.CategoryPCGames.ID
			}
		}
		attrs := []torznab.SearchResultItemTorznabAttr{
			{
				AttrName:  torznab.AttrInfoHash,
				AttrValue: item.Torrent.InfoHash.String(),
			},
			{
				AttrName:  torznab.AttrMagnetUrl,
				AttrValue: item.Torrent.MagnetUri(),
			},
			{
				AttrName:  torznab.AttrCategory,
				AttrValue: strconv.Itoa(categoryId),
			},
			{
				AttrName:  torznab.AttrSize,
				AttrValue: strconv.FormatUint(uint64(item.Torrent.Size), 10),
			},
			{
				AttrName:  torznab.AttrPublishDate,
				AttrValue: item.PublishedAt.Format(torznab.RssDateDefaultFormat),
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
		if !item.Content.ReleaseYear.IsNil() {
			attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
				AttrName:  torznab.AttrYear,
				AttrValue: strconv.Itoa(int(item.Content.ReleaseYear)),
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
			PubDate:  torznab.RssDate(item.PublishedAt),
			Enclosure: torznab.SearchResultItemEnclosure{
				URL:    item.Torrent.MagnetUri(),
				Type:   "application/x-bittorrent;x-scheme-handler/magnet",
				Length: strconv.FormatUint(uint64(item.Torrent.Size), 10),
			},
			TorznabAttrs: attrs,
		})
	}
	return torznab.SearchResult{
		Channel: torznab.SearchResultChannel{
			Title: a.title,
			Response: torznab.SearchResultResponse{
				Offset: req.Offset.Uint,
				//Total:  res.TotalCount,
			},
			Items: entries,
		},
	}
}
