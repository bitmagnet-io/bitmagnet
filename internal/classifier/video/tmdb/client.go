package tmdb

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient/httplogger"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient/httpratelimiter"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/cyruzin/golang-tmdb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
	Search search.Search
}

type Result struct {
	fx.Out
	TmdbClient *tmdb.Client
	Client     Client
}

type Client interface {
	MovieClient
	TvShowClient
}

func New(p Params) (r Result, err error) {
	logger := p.Logger.Named("tmdb_client")
	rateLimit := p.Config.RateLimit
	rateLimitBurst := p.Config.RateLimitBurst
	if p.Config.ApiKey == defaultTmdbApiKey {
		rateLimit = time.Second
		rateLimitBurst = 1
		logger.Warnln("you are using the default TMDB api key; TMDB requests will be limited to 1 per second; to remove this warning please configure a personal TMDB api key")
	}
	httpClient := http.Client{
		// need to set a non-zero value as the underlying client unfortunately sets 10 seconds as the default if none is provided;
		// this does not work well with the rate limiter; a 30 second timeout fixes this assuming a concurrency of 10 on the queue
		// (and a maximum of 2 TMDB requests per classification)
		Timeout: time.Second * 30,
		Transport: httpratelimiter.NewDecorator(
			rateLimit,
			rateLimitBurst,
		)(httplogger.NewDecorator(
			logger,
		)(http.DefaultTransport)),
	}
	c, initErr := tmdb.Init(p.Config.ApiKey)
	c.SetClientConfig(httpClient)
	if initErr != nil {
		err = initErr
		return
	}
	r.Client = &client{
		c: c,
		s: p.Search,
	}
	r.TmdbClient = c
	return
}

type client struct {
	c *tmdb.Client
	s search.Search
}

const SourceTmdb = "tmdb"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnknownSource = errors.New("unknown source")
)
