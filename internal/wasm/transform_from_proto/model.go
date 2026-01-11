package transform_from_proto

import (
	"database/sql"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	proto "github.com/bitmagnet-io/bitmagnet/proto/common/model"
)

func Torrent(pt *proto.Torrent) model.Torrent {
	infoHash, _ := protocol.ParseID(pt.InfoHash)

	return model.Torrent{
		InfoHash:    infoHash,
		Name:        pt.Name,
		Size:        uint(pt.Size),
		Private:     pt.Private,
		FilesStatus: FilesStatus(pt.FilesStatus),
		Extension:   newNullString(pt.Extension),
		FilesCount:  newNullUint(pt.FilesCount),
		Files:       slice.Map(pt.Files, TorrentFile),
		Sources:     slice.Map(pt.Sources, TorrentSource),
		Tags: slice.Map(pt.Tags, func(s string) model.TorrentTag {
			return model.TorrentTag{
				InfoHash: infoHash,
				Name:     s,
			}
		}),
		CreatedAt: time.UnixMilli(pt.CreatedAt),
		UpdatedAt: time.UnixMilli(pt.UpdatedAt),
	}
}

func TorrentFile(pf *proto.TorrentFile) model.TorrentFile {
	return model.TorrentFile{
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(pf.InfoHash)
			return infoHash
		}(),
		Index:     uint(pf.Index),
		Path:      pf.Path,
		Extension: newNullString(pf.Extension),
		Size:      uint(pf.Size),
		CreatedAt: time.UnixMilli(pf.CreatedAt),
		UpdatedAt: time.UnixMilli(pf.UpdatedAt),
	}
}

func TorrentSource(ps *proto.TorrentSource) model.TorrentsTorrentSource {
	return model.TorrentsTorrentSource{
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(ps.InfoHash)
			return infoHash
		}(),
		Source:   ps.Source,
		ImportID: newNullString(ps.ImportId),
		Seeders:  newNullUint(ps.Seeders),
		Leechers: newNullUint(ps.Leechers),
		PublishedAt: func() sql.NullTime {
			if ps.PublishedAt != nil && *ps.PublishedAt != 0 {
				return sql.NullTime{
					Time:  time.UnixMilli(*ps.PublishedAt),
					Valid: true,
				}
			}
			return sql.NullTime{}
		}(),
		CreatedAt: time.UnixMilli(ps.CreatedAt),
		UpdatedAt: time.UnixMilli(ps.UpdatedAt),
	}
}

func TorrentContent(tcp *proto.TorrentContent) model.TorrentContent {
	return model.TorrentContent{
		ID: tcp.Id,
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(tcp.InfoHash)
			return infoHash
		}(),
		ContentType: func() model.NullContentType {
			if tcp.ContentType != nil {
				ct, err := model.ParseContentType(*tcp.ContentType)
				if err == nil {
					return model.NullContentType{
						ContentType: ct,
						Valid:       true,
					}
				}
			}

			return model.NullContentType{}
		}(),
		Content: Content(tcp.Content),
		Languages: func() model.Languages {

			languages := make(model.Languages)
			for _, l := range tcp.Languages {
				lang := model.ParseLanguage(l)
				if !lang.Valid {
					continue
				}
				languages[lang.Language] = struct{}{}
			}

			if len(languages) == 0 {
				return nil
			}

			return languages
		}(),
		Episodes: func() model.Episodes {
			episodes := make(model.Episodes)
			for _, str := range tcp.Episodes {
				for season, eps := range model.ParseEpisodes(str) {
					for ep := range eps {
						episodes = episodes.AddEpisode(season, ep)
					}
				}
			}

			if len(episodes) == 0 {
				return nil
			}

			return episodes
		}(),
		VideoResolution: model.NewNullVideoResolution(tcp.VideoResolution),
		VideoSource:     model.NewNullVideoSource(tcp.VideoSource),
		Video3D:         model.NewNullVideo3D(tcp.Video3D),
		VideoModifier:   model.NewNullVideoModifier(tcp.VideoModifier),
		ReleaseGroup:    newNullString(tcp.ReleaseGroup),
		Seeders:         newNullUint(tcp.Seeders),
		Leechers:        newNullUint(tcp.Leechers),
		Size:            uint(tcp.Size),
		FilesCount:      newNullUint(tcp.FilesCount),
		Tags:            tcp.Tags,
		PublishedAt:     time.UnixMilli(tcp.PublishedAt),
		CreatedAt:       time.UnixMilli(tcp.CreatedAt),
		UpdatedAt:       time.UnixMilli(tcp.UpdatedAt),
		Torrent:         Torrent(tcp.Torrent),
	}
}

func Content(pc *proto.Content) model.Content {
	if pc == nil {
		return model.Content{}
	}

	ref := ContentRef(pc.Ref)

	return model.Content{
		Type:        ref.Type,
		Source:      ref.Source,
		ID:          ref.ID,
		Title:       pc.Title,
		ReleaseDate: Date(pc.ReleaseDate),
		Adult:       newNullBool(pc.Adult),
		OriginalLanguage: func() model.NullLanguage {
			if pc.OriginalLanguage != nil {
				lang := model.ParseLanguage(*pc.OriginalLanguage)
				if lang.Valid {
					return lang
				}
			}
			return model.NullLanguage{}
		}(),
		OriginalTitle: newNullString(pc.OriginalTitle),
		Overview:      newNullString(pc.Overview),
		Runtime:       newNullUint16(pc.Runtime),
		Popularity:    newNullFloat32(pc.Popularity),
		VoteAverage:   newNullFloat32(pc.VoteAverage),
		VoteCount:     newNullUint(pc.VoteCount),
		Tags:          slice.Map(pc.Tags, func(s string) model.ContentTag { return model.ContentTag{Name: s} }),
		Collections:   slice.Map(pc.Collections, ContentCollection),
		Attributes:    slice.Map(pc.Attributes, ContentAttribute),
		CreatedAt:     time.UnixMilli(pc.CreatedAt),
		UpdatedAt:     time.UnixMilli(pc.UpdatedAt),
	}
}

func ContentRef(r *proto.ContentRef) model.ContentRef {
	if r != nil {
		if t, err := model.ParseContentType(r.Type); err != nil {

			return model.ContentRef{
				Type:   t,
				Source: r.Source,
				ID:     r.Id,
			}
		}
	}
	return model.ContentRef{}
}

func ContentCollection(c *proto.ContentCollection) model.ContentCollection {
	return model.ContentCollection{
		Type:      c.Type,
		Source:    c.Source,
		ID:        c.Id,
		Name:      c.Name,
		CreatedAt: time.UnixMilli(c.CreatedAt),
		UpdatedAt: time.UnixMilli(c.UpdatedAt),
	}
}

func ContentAttribute(a *proto.ContentAttribute) model.ContentAttribute {
	return model.ContentAttribute{
		Source:    a.Source,
		Key:       a.Key,
		Value:     a.Value,
		CreatedAt: time.UnixMilli(a.CreatedAt),
		UpdatedAt: time.UnixMilli(a.UpdatedAt),
	}
}

func FilesStatus(fs proto.Torrent_FilesStatus) model.FilesStatus {
	switch fs {
	case proto.Torrent_single:
		return model.FilesStatusSingle
	case proto.Torrent_multi:
		return model.FilesStatusMulti
	case proto.Torrent_over_threshold:
		return model.FilesStatusOverThreshold
	default:
		return model.FilesStatusNoInfo
	}
}

func Date(d *proto.Date) model.Date {
	var result model.Date
	if d != nil {
		result.Year = model.Year(d.Year)
		if d.Month != nil && *d.Month >= int32(time.January) && *d.Month <= int32(time.December) {
			result.Month = time.Month(*d.Month)
		}
		if d.Day != nil {
			result.Day = uint8(*d.Day)
		}
	}
	return result
}

func newNullString(s *string) model.NullString {
	if s != nil {
		return model.NewNullString(*s)
	}
	return model.NullString{}
}

func newNullBool(b *bool) model.NullBool {
	if b != nil {
		return model.NewNullBool(*b)
	}
	return model.NullBool{}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func newNullUint[n number](v *n) model.NullUint {
	if v != nil {
		return model.NewNullUint(uint(*v))
	}
	return model.NullUint{}
}

func newNullUint16[n number](v *n) model.NullUint16 {
	if v != nil {
		return model.NewNullUint16(uint16(*v))
	}
	return model.NullUint16{}
}

func newNullFloat32[n number](v *n) model.NullFloat32 {
	if v != nil {
		return model.NewNullFloat32(float32(*v))
	}
	return model.NullFloat32{}
}
