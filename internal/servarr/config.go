package servarr

type UrlKey struct {
	Url    string
	ApiKey string
}

type Config struct {
	IndexerName string
	Sonarr      UrlKey
	Radarr      UrlKey
	Prowlarr    UrlKey
}

func NewDefaultConfig() Config {
	return Config{
		IndexerName: "bitmagnet",
		Sonarr:      UrlKey{Url: "http://localhost:8989", ApiKey: "private"},
		Radarr:      UrlKey{Url: "http://localhost:7878", ApiKey: "private"},
		Prowlarr:    UrlKey{Url: "http://localhost:9696", ApiKey: "private"},
	}
}
