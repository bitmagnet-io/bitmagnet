package discogs

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	discogs "github.com/irlndts/go-discogs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config      Config
	Logger      *zap.SugaredLogger
	Search      search.Search
	Persistence persistence.Persistence
}

type Result struct {
	fx.Out
	DiscogsClient discogs.Discogs
	Client        Client
}

type Client interface {
	MusicClient
}

func New(p Params) (r Result, err error) {
	logger := p.Logger.Named("discogs_client")
	if p.Config.ApiKey == defaultDiscogsApiKey {
		logger.Warn(p.Config.ApiKey)
		logger.Warnln("you are using the default Discogs api key; there is currently no default discogs api key")
		return Result{}, ErrNotApiDiscogs
	} else {

		clientDiscogs, _ := discogs.New(&discogs.Options{
			UserAgent: "MediaLibrary",
			Currency:  "EUR",                     // optional, "USD" (default), "GBP", "EUR", "CAD", "AUD", "JPY", "CHF", "MXN", "BRL", "NZD", "SEK", "ZAR" are allowed
			Token:     p.Config.ApiKey,           // optional
			URL:       "https://api.discogs.com", // optional
		})

		r.Client = &client{
			c: clientDiscogs,
			p: p.Persistence,
			s: p.Search,
		}
		r.DiscogsClient = clientDiscogs
		return
	}
}

type client struct {
	c discogs.Discogs
	s search.Search
	p persistence.Persistence
}

const SourceMDiscogs = "discogs"

var (
	ErrNotApiDiscogs = errors.New("no api key found")
	ErrNotFound      = errors.New("not found")
	ErrUnknownSource = errors.New("unknown source")
)
