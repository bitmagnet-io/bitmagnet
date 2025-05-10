package protobuf

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func NewTorrent(t model.Torrent) *Torrent {
	files := make([]*Torrent_File, 0, len(t.Files))

	var filesSize *int64

	switch t.FilesStatus {
	case model.FilesStatusNoInfo:
		ext := model.FileExtensionFromPath(t.Name)
		if ext.Valid {
			files = append(files, &Torrent_File{
				Path:      t.Name,
				BasePath:  t.Name,
				BaseName:  t.Name,
				Extension: &ext.String,
				FileType:  NewFileType(model.FileTypeFromExtension(ext.String)),
			})
			s := int64(t.Size)
			filesSize = &s
		}
	case model.FilesStatusSingle:
		var ext *string

		if t.Extension.Valid {
			str := t.Extension.String
			ext = &str
		}

		files = append(files, &Torrent_File{
			Path:      t.Name,
			BasePath:  t.BaseName(),
			BaseName:  t.BaseName(),
			Extension: ext,
			Size:      int64(t.Size),
			FileType:  NewFileType(t.FileType()),
		})
		s := int64(t.Size)
		filesSize = &s
	case model.FilesStatusMulti, model.FilesStatusOverThreshold:
		fs := int64(0)
		for _, f := range t.Files {
			fs += int64(f.Size)

			var ext *string
			if f.Extension.Valid {
				ext = &f.Extension.String
			}

			files = append(files, &Torrent_File{
				Index:     int32(f.Index),
				Path:      f.Path,
				BasePath:  f.BasePath(),
				BaseName:  f.BaseName(),
				Size:      int64(f.Size),
				Extension: ext,
				FileType:  NewFileType(f.FileType()),
			})
		}

		if len(files) == 0 && t.FilesStatus == model.FilesStatusOverThreshold {
			s := int64(t.Size)
			filesSize = &s
		} else {
			filesSize = &fs
		}
	}

	var filesCount *int32

	if t.FilesCount.Valid {
		c := int32(t.FilesCount.Uint)
		filesCount = &c
	}

	var seeders *int32

	var leechers *int32

	if s := t.Seeders(); s.Valid {
		s := int32(s.Uint)
		seeders = &s
	}

	if l := t.Leechers(); l.Valid {
		l := int32(l.Uint)
		leechers = &l
	}

	return &Torrent{
		InfoHash:           t.InfoHash.String(),
		Name:               t.Name,
		BaseName:           t.BaseName(),
		Size:               int64(t.Size),
		Files:              files,
		FilesCount:         filesCount,
		FilesSize:          filesSize,
		FileExtensions:     t.FileExtensions(),
		Seeders:            seeders,
		Leechers:           leechers,
		HasHint:            !t.Hint.IsNil(),
		HasHintedContentId: !t.Hint.IsNil() && t.Hint.ContentSource.Valid,
	}
}

func NewClassification(c classification.Result) *Classification {
	// var year *int32
	// if !c.Year.IsNil() {
	// 	y := int32(c.Year)
	// 	year = &y
	// }
	languages := slice.Map(c.Languages.Slice(), func(l model.Language) string {
		return l.ID()
	})

	episodes := slice.Map(c.Episodes.SeasonEntries(), func(e model.Season) string {
		return e.String()
	})

	var videoResolution *string

	if c.VideoResolution.Valid {
		str := c.VideoResolution.VideoResolution.String()
		videoResolution = &str
	}

	var videoSource *string

	if c.VideoSource.Valid {
		str := c.VideoSource.VideoSource.String()
		videoSource = &str
	}

	var videoCodec *string

	if c.VideoCodec.Valid {
		str := c.VideoCodec.VideoCodec.String()
		videoCodec = &str
	}

	var releaseGroup *string
	if c.ReleaseGroup.Valid {
		releaseGroup = &c.ReleaseGroup.String
	}

	var contentID *string

	var contentSource *string

	if c.Content != nil {
		contentID = &c.Content.ID
		contentSource = &c.Content.Source
	}

	return &Classification{
		ContentType:        NewContentType(c.ContentType),
		HasAttachedContent: c.Content != nil,
		HasBaseTitle:       c.BaseTitle.Valid,
		// Year:               year,
		Languages:       languages,
		Episodes:        episodes,
		VideoResolution: videoResolution,
		VideoSource:     videoSource,
		VideoCodec:      videoCodec,
		ReleaseGroup:    releaseGroup,
		ContentId:       contentID,
		ContentSource:   contentSource,
	}
}

func NewContentType(ct model.NullContentType) Classification_ContentType {
	if ct.Valid {
		switch ct.ContentType {
		case model.ContentTypeMovie:
			return Classification_movie
		case model.ContentTypeTvShow:
			return Classification_tv_show
		case model.ContentTypeMusic:
			return Classification_music
		case model.ContentTypeEbook:
			return Classification_ebook
		case model.ContentTypeComic:
			return Classification_comic
		case model.ContentTypeAudiobook:
			return Classification_audiobook
		case model.ContentTypeGame:
			return Classification_game
		case model.ContentTypeSoftware:
			return Classification_software
		case model.ContentTypeXxx:
			return Classification_xxx
		}
	}

	return Classification_unknown
}

func NewFileType(ft model.NullFileType) Torrent_File_FileType {
	if ft.Valid {
		switch ft.FileType {
		case model.FileTypeArchive:
			return Torrent_File_archive
		case model.FileTypeAudio:
			return Torrent_File_audio
		case model.FileTypeData:
			return Torrent_File_data
		case model.FileTypeImage:
			return Torrent_File_image
		case model.FileTypeSoftware:
			return Torrent_File_software
		case model.FileTypeSubtitles:
			return Torrent_File_subtitles
		case model.FileTypeVideo:
			return Torrent_File_video
		}
	}

	return Torrent_File_unknown
}
