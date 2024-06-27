package servarr

type IndexerResource struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
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
