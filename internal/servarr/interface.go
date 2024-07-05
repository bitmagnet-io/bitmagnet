package servarr

type IndexerBulkResource struct {
	Ids                     []*int64 `json:"ids,omitempty"`
	EnableRss               *bool    `json:"enableRss,omitempty"`
	EnableAutomaticSearch   *bool    `json:"enableAutomaticSearch,omitempty"`
	EnableInteractiveSearch *bool    `json:"enableInteractiveSearch,omitempty"`
}

type IndexerResource struct {
	ID                      int64  `json:"id"`
	Name                    string `json:"name,omitempty"`
	EnableRss               bool   `json:"enableRss"`
	EnableAutomaticSearch   bool   `json:"enableAutomaticSearch"`
	EnableInteractiveSearch bool   `json:"enableInteractiveSearch"`
}

type SeriesResource struct {
	ID    int64  `json:"id"`
	Title string `json:"title,omitempty"`
}

type EpisodeResource struct {
	ID            int64  `json:"id"`
	SeriesID      int64  `json:"seriesId"`
	TvdbID        int64  `json:"tvdbId"`
	EpisodeFileID int64  `json:"episodeFileId"`
	SeasonNumber  int64  `json:"seasonNumber"`
	EpisodeNumber int64  `json:"episodeNumber"`
	Title         string `json:"title,omitempty"`
}

type MovieResource struct {
	ID    int64  `json:"id"`
	Title string `json:"title,omitempty"`
}

type ReleaseResource struct {
	ID        int64  `json:"id"`
	GUID      string `json:"guid,omitempty"`
	IndexerID int64  `json:"indexerId"`
	InfoHash  string `json:"infoHash,omitempty"`
}
