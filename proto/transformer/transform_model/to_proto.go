package transform_model

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	proto "github.com/bitmagnet-io/bitmagnet/proto/common/model"
)

func TorrentToProto(t model.Torrent) *proto.Torrent {
	pt := &proto.Torrent{
		InfoHash:    t.InfoHash.String(),
		Name:        t.Name,
		Size:        int64(t.Size),
		Private:     t.Private,
		FilesStatus: FilesStatusToProto(t.FilesStatus),
		Extension: func() *string {
			if t.Extension.Valid {
				return &t.Extension.String
			}

			return nil
		}(),
		FilesCount: func() *int32 {
			if t.FilesCount.Valid {
				value := int32(t.FilesCount.Uint)

				return &value
			}

			return nil
		}(),
		Files:   slice.Map(t.Files, TorrentFileToProto),
		Sources: slice.Map(t.Sources, TorrentSourceToProto),
		Tags: slice.Map(t.Tags, func(t model.TorrentTag) string {
			return t.Name
		}),
		CreatedAt: t.CreatedAt.UnixMilli(),
		UpdatedAt: t.UpdatedAt.UnixMilli(),
	}

	pt.EnsureTimes(time.Now().UnixMilli())

	return pt
}

func TorrentFileToProto(f model.TorrentFile) *proto.TorrentFile {
	return &proto.TorrentFile{
		InfoHash: f.InfoHash.String(),
		Index:    int32(f.Index),
		Path:     f.Path,
		Extension: func() *string {
			if f.Extension.Valid {
				return &f.Extension.String
			}

			return nil
		}(),
		FileType: func() *string {
			if ft := f.FileType(); ft.Valid {
				str := ft.FileType.String()
				return &str
			}

			return nil
		}(),
		Size:      int64(f.Size),
		CreatedAt: f.CreatedAt.UnixMilli(),
		UpdatedAt: f.UpdatedAt.UnixMilli(),
	}
}

func TorrentSourceToProto(s model.TorrentsTorrentSource) *proto.TorrentSource {
	return &proto.TorrentSource{
		InfoHash: s.InfoHash.String(),
		Source:   s.Source,
		ImportId: func() *string {
			if s.ImportID.Valid {
				return &s.ImportID.String
			}

			return nil
		}(),
		Seeders: func() *int32 {
			if s.Seeders.Valid {
				value := int32(s.Seeders.Uint)

				return &value
			}

			return nil
		}(),
		Leechers: func() *int32 {
			if s.Leechers.Valid {
				value := int32(s.Leechers.Uint)

				return &value
			}

			return nil
		}(),
		PublishedAt: func() *int64 {
			if s.PublishedAt.Valid {
				value := s.PublishedAt.Time.UnixMilli()

				return &value
			}

			return nil
		}(),
		CreatedAt: s.CreatedAt.UnixMilli(),
		UpdatedAt: s.UpdatedAt.UnixMilli(),
	}
}

func TorrentContentToProto(tc model.TorrentContent) *proto.TorrentContent {
	tcp := &proto.TorrentContent{
		Id:       tc.InferID(),
		InfoHash: tc.InfoHash.String(),
		ContentType: func() *string {
			if tc.ContentType.Valid {
				value := tc.ContentType.ContentType.String()

				return &value
			}

			return nil
		}(),
		ContentRef: func() *proto.ContentRef {
			if tc.Content.ID != "" {
				return ContentRefToProto(tc.Content.Ref())
			}

			return nil
		}(),
		Languages: slice.Map(tc.Languages.Slice(), func(l model.Language) string {
			return l.ID()
		}),
		Episodes: slice.Map(tc.Episodes.SeasonEntries(), func(s model.Season) string {
			return s.String()
		}),
		VideoResolution: func() *string {
			if tc.VideoResolution.Valid {
				value := tc.VideoResolution.VideoResolution.String()

				return &value
			}

			return nil
		}(),
		VideoSource: func() *string {
			if tc.VideoSource.Valid {
				value := tc.VideoSource.VideoSource.String()

				return &value
			}

			return nil
		}(),
		Video3D: func() *string {
			if tc.Video3D.Valid {
				value := tc.Video3D.Video3D.String()

				return &value
			}

			return nil
		}(),
		VideoModifier: func() *string {
			if tc.VideoModifier.Valid {
				value := tc.VideoModifier.VideoModifier.String()

				return &value
			}

			return nil
		}(),
		ReleaseGroup: func() *string {
			if tc.ReleaseGroup.Valid {
				return &tc.ReleaseGroup.String
			}

			return nil
		}(),
		Seeders: func() *int32 {
			if tc.Seeders.Valid {
				value := int32(tc.Seeders.Uint)

				return &value
			}

			return nil
		}(),
		Leechers: func() *int32 {
			if tc.Leechers.Valid {
				value := int32(tc.Leechers.Uint)

				return &value
			}

			return nil
		}(),
		Size: int64(tc.Size),
		FilesCount: func() *int32 {
			if tc.FilesCount.Valid {
				value := int32(tc.FilesCount.Uint)
				return &value
			}

			return nil
		}(),
		Tags:        tc.Tags,
		PublishedAt: tc.PublishedAt.UnixMilli(),
		CreatedAt:   tc.CreatedAt.UnixMilli(),
		UpdatedAt:   tc.UpdatedAt.UnixMilli(),
		Torrent:     TorrentToProto(tc.Torrent),
		Content: func() *proto.Content {
			if tc.Content.ID != "" {
				return ContentToProto(tc.Content)
			}

			return nil
		}(),
	}

	tcp.EnsureTimes(time.Now().UnixMilli())

	return tcp
}

func ContentToProto(c model.Content) *proto.Content {
	cp := &proto.Content{
		Ref:         ContentRefToProto(c.Ref()),
		Title:       c.Title,
		ReleaseDate: DateToProto(c.ReleaseDate),
		Adult: func() *bool {
			if c.Adult.Valid {
				return &c.Adult.Bool
			}

			return nil
		}(),
		OriginalLanguage: func() *string {
			if c.OriginalLanguage.Valid {
				value := c.OriginalLanguage.Language.String()

				return &value
			}

			return nil
		}(),
		OriginalTitle: func() *string {
			if c.OriginalTitle.Valid {
				return &c.OriginalTitle.String
			}

			return nil
		}(),
		Overview: func() *string {
			if c.Overview.Valid {
				return &c.Overview.String
			}

			return nil
		}(),
		Runtime: func() *int32 {
			if c.Runtime.Valid {
				value := int32(c.Runtime.Uint16)

				return &value
			}

			return nil
		}(),
		Popularity: func() *float32 {
			if c.Popularity.Valid {
				return &c.Popularity.Float32
			}

			return nil
		}(),
		VoteAverage: func() *float32 {
			if c.VoteAverage.Valid {
				return &c.VoteAverage.Float32
			}

			return nil
		}(),
		VoteCount: func() *int32 {
			if c.VoteCount.Valid {
				value := int32(c.VoteCount.Uint)

				return &value
			}

			return nil
		}(),
		Tags: slice.Map(c.Tags, func(t model.ContentTag) string {
			return t.Name
		}),
		Collections: slice.Map(c.Collections, ContentCollectionToProto),
		Attributes:  slice.Map(c.Attributes, ContentAttributeToProto),
		CreatedAt:   c.CreatedAt.UnixMilli(),
		UpdatedAt:   c.UpdatedAt.UnixMilli(),
	}

	cp.EnsureTimes(time.Now().UnixMilli())

	return cp
}

func ContentRefToProto(r model.ContentRef) *proto.ContentRef {
	return &proto.ContentRef{
		Type:   r.Type.String(),
		Source: r.Source,
		Id:     r.ID,
	}
}

func ContentCollectionToProto(c model.ContentCollection) *proto.ContentCollection {
	return &proto.ContentCollection{
		Type:      c.Type,
		Source:    c.Source,
		Id:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.UnixMilli(),
		UpdatedAt: c.UpdatedAt.UnixMilli(),
	}
}

func ContentAttributeToProto(a model.ContentAttribute) *proto.ContentAttribute {
	return &proto.ContentAttribute{
		Source:    a.Source,
		Key:       a.Key,
		Value:     a.Value,
		CreatedAt: a.CreatedAt.UnixMilli(),
		UpdatedAt: a.UpdatedAt.UnixMilli(),
	}
}

func FilesStatusToProto(fs model.FilesStatus) proto.Torrent_FilesStatus {
	switch fs {
	case model.FilesStatusSingle:
		return proto.Torrent_single
	case model.FilesStatusMulti:
		return proto.Torrent_multi
	case model.FilesStatusOverThreshold:
		return proto.Torrent_over_threshold
	default:
		return proto.Torrent_no_info
	}
}

func DateToProto(d model.Date) *proto.Date {
	if d.IsNil() {
		return nil
	}

	return &proto.Date{
		Year: int32(d.Year),
		Month: func() *int32 {
			if d.Month > 0 {
				value := int32(d.Month)

				return &value
			}

			return nil
		}(),
		Day: func() *int32 {
			if d.Day > 0 {
				value := int32(d.Day)

				return &value
			}

			return nil
		}(),
	}
}
