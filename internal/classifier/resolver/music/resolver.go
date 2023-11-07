package music

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/music/discogs"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"strings"
)

type Params struct {
	fx.In
	DiscogsClient discogs.Client
}

type Result struct {
	fx.Out
	Resolver resolver.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: musicResolver{
			config:        resolver.SubResolverConfig{Key: "music", Priority: 3},
			discogsClient: p.DiscogsClient,
		},
	}
}

type musicResolver struct {
	config        resolver.SubResolverConfig
	discogsClient discogs.Client
}

func (r musicResolver) Config() resolver.SubResolverConfig {
	return r.config
}

func (r musicResolver) PreEnrich(content model.TorrentContent) (model.TorrentContent, error) {
	return PreEnrich(content)
}

func (r musicResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	titleLower := strings.ToLower(content.Torrent.Name)

	fmt.Printf("Try to resolve : %s\n", titleLower)

	if strings.Contains(titleLower, "discography") ||
		strings.Contains(titleLower, "discographie") ||
		strings.Contains(titleLower, "discografia") ||
		strings.Contains(titleLower, "anthology") {

		artist, err := FindArtistDiscography(titleLower)
		if err != nil {
			fmt.Printf("Pas trouvé artist : %s\n", content.Torrent.Name)
			return content, err
		}

		resultArtist, err := r.discogsClient.SearchArtist(ctx, discogs.SearchMusicParams{
			Artist:               artist,
			Album:                "",
			Track:                "",
			LevenshteinThreshold: 0,
		})
		fmt.Printf("Discography : %s\n", artist)

		content.ContentType.Valid = true
		content.ContentType.ContentType = model.ContentTypeMusic
		content.Content = resultArtist
		return content, err
	} else {
		for _, ext := range model.FileTypeAudio.Extensions() {
			if strings.Contains(titleLower, ext) {
				fmt.Printf("We may hava found musics %s\n", content.Torrent.Name)
				content.ContentType.Valid = true
				content.ContentType.ContentType = model.ContentTypeMusic

				return content, nil
			}
		}
	}
	return content, nil
}
