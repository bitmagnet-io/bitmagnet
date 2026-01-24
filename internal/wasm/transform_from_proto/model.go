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
	infoHash, _ := protocol.ParseID(pt.GetInfoHash())

	return model.Torrent{
		InfoHash:    infoHash,
		Name:        pt.GetName(),
		Size:        uint(pt.GetSize()),
		Private:     pt.GetPrivate(),
		FilesStatus: FilesStatus(pt.GetFilesStatus()),
		Extension:   newNullString(pt.GetExtension()),
		FilesCount:  newNullUint(pt.GetFilesCount()),
		Files:       slice.Map(pt.GetFiles(), TorrentFile),
		Sources:     slice.Map(pt.GetSources(), TorrentSource),
		Tags: slice.Map(pt.GetTags(), func(s string) model.TorrentTag {
			return model.TorrentTag{
				InfoHash: infoHash,
				Name:     s,
			}
		}),
		CreatedAt: time.UnixMilli(pt.GetCreatedAt()),
		UpdatedAt: time.UnixMilli(pt.GetUpdatedAt()),
	}
}

func TorrentFile(pf *proto.TorrentFile) model.TorrentFile {
	return model.TorrentFile{
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(pf.GetInfoHash())
			return infoHash
		}(),
		Index:     uint(pf.GetIndex()),
		Path:      pf.GetPath(),
		Extension: newNullString(pf.GetExtension()),
		Size:      uint(pf.GetSize()),
		CreatedAt: time.UnixMilli(pf.GetCreatedAt()),
		UpdatedAt: time.UnixMilli(pf.GetUpdatedAt()),
	}
}

func TorrentSource(ps *proto.TorrentSource) model.TorrentsTorrentSource {
	return model.TorrentsTorrentSource{
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(ps.GetInfoHash())
			return infoHash
		}(),
		Source:   ps.GetSource(),
		ImportID: newNullString(ps.GetImportId()),
		Seeders:  newNullUint(ps.GetSeeders()),
		Leechers: newNullUint(ps.GetLeechers()),
		PublishedAt: func() sql.NullTime {
			if publishedAt := ps.GetPublishedAt(); publishedAt != nil && *publishedAt != 0 {
				return sql.NullTime{
					Time:  time.UnixMilli(*publishedAt),
					Valid: true,
				}
			}

			return sql.NullTime{}
		}(),
		CreatedAt: time.UnixMilli(ps.GetCreatedAt()),
		UpdatedAt: time.UnixMilli(ps.GetUpdatedAt()),
	}
}

func TorrentContent(tcp *proto.TorrentContent) model.TorrentContent {
	return model.TorrentContent{
		ID: tcp.GetId(),
		InfoHash: func() protocol.ID {
			infoHash, _ := protocol.ParseID(tcp.GetInfoHash())
			return infoHash
		}(),
		ContentType: func() model.NullContentType {
			if ct := tcp.GetContentType(); ct != nil {
				ct, err := model.ParseContentType(*ct)
				if err == nil {
					return model.NullContentType{
						ContentType: ct,
						Valid:       true,
					}
				}
			}

			return model.NullContentType{}
		}(),
		Content: Content(tcp.GetContent()),
		Languages: func() model.Languages {
			languages := make(model.Languages)

			for _, l := range tcp.GetLanguages() {
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

			for _, str := range tcp.GetEpisodes() {
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
		VideoResolution: model.NewNullVideoResolution(tcp.GetVideoResolution()),
		VideoSource:     model.NewNullVideoSource(tcp.GetVideoSource()),
		Video3D:         model.NewNullVideo3D(tcp.GetVideo3D()),
		VideoModifier:   model.NewNullVideoModifier(tcp.GetVideoModifier()),
		ReleaseGroup:    newNullString(tcp.GetReleaseGroup()),
		Seeders:         newNullUint(tcp.GetSeeders()),
		Leechers:        newNullUint(tcp.GetLeechers()),
		Size:            uint(tcp.GetSize()),
		FilesCount:      newNullUint(tcp.GetFilesCount()),
		Tags:            tcp.GetTags(),
		PublishedAt:     time.UnixMilli(tcp.GetPublishedAt()),
		CreatedAt:       time.UnixMilli(tcp.GetCreatedAt()),
		UpdatedAt:       time.UnixMilli(tcp.GetUpdatedAt()),
		Torrent:         Torrent(tcp.GetTorrent()),
	}
}

func Content(pc *proto.Content) model.Content {
	if pc == nil {
		return model.Content{}
	}

	ref := ContentRef(pc.GetRef())

	return model.Content{
		Type:        ref.Type,
		Source:      ref.Source,
		ID:          ref.ID,
		Title:       pc.GetTitle(),
		ReleaseDate: Date(pc.GetReleaseDate()),
		Adult:       newNullBool(pc.GetAdult()),
		OriginalLanguage: func() model.NullLanguage {
			if l := pc.GetOriginalLanguage(); l != nil {
				lang := model.ParseLanguage(*l)
				if lang.Valid {
					return lang
				}
			}

			return model.NullLanguage{}
		}(),
		OriginalTitle: newNullString(pc.GetOriginalTitle()),
		Overview:      newNullString(pc.GetOverview()),
		Runtime:       newNullUint16(pc.GetRuntime()),
		Popularity:    newNullFloat32(pc.GetPopularity()),
		VoteAverage:   newNullFloat32(pc.GetVoteAverage()),
		VoteCount:     newNullUint(pc.GetVoteCount()),
		Tags: slice.Map(
			pc.GetTags(),
			func(s string) model.ContentTag { return model.ContentTag{Name: s} },
		),
		Collections: slice.Map(pc.GetCollections(), ContentCollection),
		Attributes:  slice.Map(pc.GetAttributes(), ContentAttribute),
		CreatedAt:   time.UnixMilli(pc.GetCreatedAt()),
		UpdatedAt:   time.UnixMilli(pc.GetUpdatedAt()),
	}
}

func ContentRef(r *proto.ContentRef) model.ContentRef {
	if r != nil {
		if t, err := model.ParseContentType(r.GetType()); err != nil {
			return model.ContentRef{
				Type:   t,
				Source: r.GetSource(),
				ID:     r.GetId(),
			}
		}
	}

	return model.ContentRef{}
}

func ContentCollection(c *proto.ContentCollection) model.ContentCollection {
	return model.ContentCollection{
		Type:      c.GetType(),
		Source:    c.GetSource(),
		ID:        c.GetId(),
		Name:      c.GetName(),
		CreatedAt: time.UnixMilli(c.GetCreatedAt()),
		UpdatedAt: time.UnixMilli(c.GetUpdatedAt()),
	}
}

func ContentAttribute(a *proto.ContentAttribute) model.ContentAttribute {
	return model.ContentAttribute{
		Source:    a.GetSource(),
		Key:       a.GetKey(),
		Value:     a.GetValue(),
		CreatedAt: time.UnixMilli(a.GetCreatedAt()),
		UpdatedAt: time.UnixMilli(a.GetUpdatedAt()),
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
		result.Year = model.Year(d.GetYear())
		if month := d.GetMonth(); month != nil && *month >= int32(time.January) &&
			*month <= int32(time.December) {
			result.Month = time.Month(*month)
		}

		if day := d.GetDay(); day != nil {
			result.Day = uint8(*day)
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
