package tpdb

import (
	"context"
	"fmt"
	"strconv"

	porndb "git.sr.ht/~dragnel/go-tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type MovieClient interface {
	SearchScene(ctx context.Context, p string) (model.Content, error)
	// GetMovieByExternalId(ctx context.Context, source, id string) (model.Content, error)
}

type SearchMovieParams struct {
	Title string
}

func (c *client) SearchScene(ctx context.Context, p string) (movie model.Content, err error) {
	scene, err := c.c.Parse(p)
	if err != nil {
		return model.Content{}, err
	}

	fmt.Printf("TPDB: found : %s\n", scene.Title)
	return SceneToXxxModel(scene)
}

func SceneToXxxModel(scene porndb.Scene) (movie model.Content, err error) {
	releaseDate := model.Date{}
	if scene.Date != "" {
		parsedDate, parseDateErr := model.NewDateFromIsoString(scene.Date)
		if parseDateErr != nil {
			err = parseDateErr
			return
		}
		releaseDate = parsedDate
	}
	searchString := scene.Title + " " + scene.Site.Name
	var collections []model.ContentCollection
	for _, tag := range scene.Tags {
		collections = append(collections, model.ContentCollection{
			Type:   "tag",
			Source: SourceTpdb,
			ID:     strconv.Itoa(int(tag.Id)),
			Name:   tag.Name,
		})
		searchString += " " + tag.Name
	}

	for _, performer := range scene.Performers {
		collections = append(collections, model.ContentCollection{
			Type:   "performer",
			Source: SourceTpdb,
			ID:     performer.Id,
			Name:   performer.Name,
		})
		searchString += " " + performer.Name
	}
	var attributes []model.ContentAttribute

	attributes = append(attributes, model.ContentAttribute{
		Source: SourceTpdb,
		Key:    "site",
		Value:  scene.Site.Name,
	})
	attributes = append(attributes, model.ContentAttribute{
		Source: SourceTpdb,
		Key:    "image",
		Value:  scene.Image,
	})
	attributes = append(attributes, model.ContentAttribute{
		Source: SourceTpdb,
		Key:    "poster",
		Value:  scene.Poster,
	})
	attributes = append(attributes, model.ContentAttribute{
		Source: SourceTpdb,
		Key:    "bg",
		Value:  scene.Background.Small,
	})
	releaseYear := releaseDate.Year

	return model.Content{
		Type:          model.ContentTypeXxx,
		Source:        SourceTpdb,
		ID:            scene.Id,
		Title:         scene.Title,
		ReleaseDate:   releaseDate,
		ReleaseYear:   releaseYear,
		Adult:         model.NewNullBool(true),
		OriginalTitle: model.NewNullString(scene.Title),
		Overview: model.NullString{
			String: scene.Description,
			Valid:  true,
		},
		Runtime: model.NullUint16{
			Uint16: uint16(scene.Duration),
			Valid:  true,
		},
		SearchString: searchString,
		Collections:  collections,
		Attributes:   attributes,
	}, nil
}
