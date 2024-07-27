package client

type Config struct {
	ArrServiceUrl   string
	DownloadEnabled bool
	DownloadClient  string
}

func NewDefaultConfig() Config {
	return Config{
		ArrServiceUrl:   "http://localhost:3335",
		DownloadEnabled: true,
		DownloadClient:  "Servarr",
	}
}
