package adapter

import (
  "strconv"

  "github.com/bitmagnet-io/bitmagnet/internal/database/search"
  "github.com/bitmagnet-io/bitmagnet/internal/model"
  "github.com/bitmagnet-io/bitmagnet/internal/torznab"
)

func torrentContentResultToTorznabResult(
  req torznab.SearchRequest,
  res search.TorrentContentResult,
) torznab.SearchResult {
  entries := make([]torznab.SearchResultItem, 0, len(res.Items))
  for _, item := range res.Items {
    entries = append(entries, torrentContentResultItemToTorznabResultItem(item))
  }
  return torznab.SearchResult{
    Channel: torznab.SearchResultChannel{
      Title: req.Profile.Title,
      Response: torznab.SearchResultResponse{
        Offset: req.Offset.Uint,
        //Total:  res.TotalCount,
      },
      Items: entries,
    },
  }
}

func torrentContentResultItemToTorznabResultItem(item search.TorrentContentResultItem) torznab.SearchResultItem {
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
      AttrName:  torznab.AttrMagnetURL,
      AttrValue: item.Torrent.MagnetURI(),
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
  seeders := item.Torrent.Seeders()
  leechers := item.Torrent.Leechers()
  if seeders.Valid {
    attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
      AttrName:  torznab.AttrSeeders,
      AttrValue: strconv.Itoa(int(seeders.Uint)),
    })
  }
  if leechers.Valid {
    attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
      AttrName:  torznab.AttrLeechers,
      AttrValue: strconv.Itoa(int(leechers.Uint)),
    })
  }
  if leechers.Valid && seeders.Valid {
    attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
      AttrName:  torznab.AttrPeers,
      AttrValue: strconv.Itoa(int(leechers.Uint) + int(seeders.Uint)),
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
  if tmdbid, ok := item.Content.Identifier("tmdb"); ok {
    attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
      AttrName:  torznab.AttrTmdb,
      AttrValue: tmdbid,
    })
  }
  if imdbId, ok := item.Content.Identifier("imdb"); ok {
    attrs = append(attrs, torznab.SearchResultItemTorznabAttr{
      AttrName:  torznab.AttrImdb,
      AttrValue: imdbId[2:],
    })
  }
  return torznab.SearchResultItem{
    Title:    item.Torrent.Name,
    Size:     item.Torrent.Size,
    Category: category,
    GUID:     item.InfoHash.String(),
    PubDate:  torznab.RSSDate(item.PublishedAt),
    Enclosure: torznab.SearchResultItemEnclosure{
      URL:    item.Torrent.MagnetURI(),
      Type:   "application/x-bittorrent;x-scheme-handler/magnet",
      Length: strconv.FormatUint(uint64(item.Torrent.Size), 10),
    },
    TorznabAttrs: attrs,
  }
}
