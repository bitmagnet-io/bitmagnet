package importer

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type Torrent struct {
	Source          string                    `json:"source"`
	InfoHash        protocol.ID               `json:"infoHash"`
	Name            string                    `json:"name"`
	Size            uint                      `json:"size"`
	Private         bool                      `json:"private"`
	ContentType     model.NullContentType     `json:"contentType"`
	ContentSource   model.NullString          `json:"contentSource"`
	ContentID       model.NullString          `json:"contentId"`
	Title           model.NullString          `json:"title"`
	ReleaseDate     model.Date                `json:"releaseDate"`
	ReleaseYear     model.Year                `json:"releaseYear"`
	Episodes        model.Episodes            `json:"episodes"`
	VideoResolution model.NullVideoResolution `json:"videoResolution"`
	VideoSource     model.NullVideoSource     `json:"videoSource"`
	VideoCodec      model.NullVideoCodec      `json:"videoCodec"`
	Video3D         model.NullVideo3D         `json:"video3D"`
	VideoModifier   model.NullVideoModifier   `json:"videoModifier"`
	ReleaseGroup    model.NullString          `json:"releaseGroup"`
	PublishedAt     time.Time                 `json:"publishedAt"`
	FilesCount      model.NullUint            `json:"filesCount"`
	Files           []TorrentFile             `json:"files"`
}

type TorrentFile struct {
	Index model.NullUint `json:"index"`
	Path  string         `json:"path"`
	Size  uint           `json:"size"`
}
