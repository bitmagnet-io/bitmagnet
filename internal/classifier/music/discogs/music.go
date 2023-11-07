package discogs

import (
	"context"
	"fmt"
	// "errors"
	// "github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"strconv"

	// "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	// "github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/irlndts/go-discogs"
)

type MusicClient interface {
	SearchMusic(ctx context.Context, p SearchMusicParams) (model.Content, error)
	SearchArtist(ctx context.Context, p SearchMusicParams) (model.Content, error)
}

type SearchMusicParams struct {
	Artist               string
	Album                string
	Track                string
	LevenshteinThreshold uint
}

func (c *client) SearchArtist(ctx context.Context, p SearchMusicParams) (model.Content, error) {
	return c.searchArtistDiscogs(ctx, p)
}

func (c *client) SearchMusic(ctx context.Context, p SearchMusicParams) (movie model.Content, err error) {
	fmt.Printf("SearchMusic : %s\n", p.Artist)
	_, err = c.searchArtistDiscogs(ctx, p)
	err = ErrNotFound
	return movie, err
}

func (c *client) searchMusicLocal(ctx context.Context, p SearchMusicParams) (movie model.Content, err error) {
	return model.Content{}, ErrNotFound
}

func (c *client) searchArtistDiscogs(ctx context.Context, p SearchMusicParams) (movie model.Content, err error) {
	request := discogs.SearchRequest{Artist: p.Artist, Page: 0, PerPage: 1}
	resultsDiscogs, _ := c.c.Search(request)

	var collections []model.ContentCollection

	if len(resultsDiscogs.Results) > 0 {

		resultsGenre := resultsDiscogs.Results[0].Genre
		for _, genre := range resultsGenre {
			fmt.Println(genre)
			request := discogs.SearchRequest{Genre: genre, Page: 0, PerPage: 1}
			resultsDiscogs, _ := c.c.Search(request)
			resultGenre := resultsDiscogs.Results[0]

			collections = append(collections, model.ContentCollection{
				Type:   "genre",
				Source: SourceDiscogs,
				ID:     strconv.Itoa(resultGenre.ID),
				Name:   genre,
			})
		}

		return model.Content{
			Type:             model.ContentTypeMusic,
			Source:           SourceDiscogs,
			ID:               strconv.Itoa(int(resultsDiscogs.Results[0].ID)),
			Title:            p.Artist + " Discography",
			ReleaseDate:      model.Date{},
			Adult:            model.NewNullBool(false),
			OriginalLanguage: model.NullLanguage{},
			OriginalTitle:    model.NullString{},
			Overview:         model.NullString{},
			Runtime:          model.NullUint16{},
			Popularity:       model.NullFloat32{},
			VoteAverage:      model.NullFloat32{},
			VoteCount:        model.NullUint{},
			Collections:      collections,
			Attributes:       []model.ContentAttribute{},
		}, nil
	}

	return model.Content{}, err
}

var SourceDiscogs = "discogs"
