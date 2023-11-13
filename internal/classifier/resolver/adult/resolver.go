package adult

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult/tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"strings"
)

type Params struct {
	fx.In
	TpdbClient tpdb.Client
}

type Result struct {
	fx.Out
	Resolver resolver.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: adultResolver{
			config:     resolver.SubResolverConfig{Key: "adult", Priority: 2},
			tpdbClient: p.TpdbClient,
		},
	}
}

type adultResolver struct {
	config     resolver.SubResolverConfig
	tpdbClient tpdb.Client
}

func (r adultResolver) Config() resolver.SubResolverConfig {
	return r.config
}

func (r adultResolver) PreEnrich(content model.TorrentContent) (model.TorrentContent, error) {
	//return content, nil
	return PreEnrich(content)
}

func (r adultResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {

	if strings.Contains(strings.ToLower(content.Title), "xxx") {

		if r.tpdbClient != nil {
			contentAdult, err := r.tpdbClient.SearchScene(ctx, content.Title)
			if err == nil {
				content.Title = contentAdult.Title
				content.ContentType.Valid = true
				content.Content = contentAdult
				content.SearchString = contentAdult.SearchString
				return content, nil
			}
			return model.TorrentContent{}, resolver.ErrNoMatch
		}
		content.ContentType.Valid = true
		content.ContentType.ContentType = model.ContentTypeXxx
		return content, nil
	}
	return model.TorrentContent{}, resolver.ErrNoMatch
}
