package protobuf

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
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
				Extension: &ext.String,
				FileType:  NewFileType(model.FileTypeFromExtension(ext.String)),
			})
			s := int64(t.Size)
			filesSize = &s
		}
	case model.FilesStatusSingle:
		files = append(files, &Torrent_File{
			Path:     t.Name,
			Size:     int64(t.Size),
			FileType: NewFileType(t.FileType()),
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
				Path:      f.Path,
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
		InfoHash:       t.InfoHash.String(),
		Name:           t.Name,
		Size:           int64(t.Size),
		Files:          files,
		FilesCount:     filesCount,
		FilesSize:      filesSize,
		FileExtensions: t.FileExtensions(),
		Seeders:        seeders,
		Leechers:       leechers,
	}
}

func NewClassification(c classifier.Classification) *Classification {
	var year *int32
	if !c.Year.IsNil() {
		y := int32(c.Year)
		year = &y
	}
	var languages []string
	for _, l := range c.Languages.Slice() {
		languages = append(languages, l.Id())
	}
	var episodes []string
	for _, e := range c.Episodes.SeasonEntries() {
		episodes = append(episodes, e.String())
	}
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
	return &Classification{
		ContentType:     NewContentType(c.ContentType),
		Year:            year,
		Languages:       languages,
		Episodes:        episodes,
		VideoResolution: videoResolution,
		VideoSource:     videoSource,
		VideoCodec:      videoCodec,
		ReleaseGroup:    releaseGroup,
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
