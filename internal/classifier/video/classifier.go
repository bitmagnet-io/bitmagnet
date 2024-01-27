package video

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type videoClassifier struct {
	tmdbClient tmdb.Client
}

func (c videoClassifier) Key() string {
	return "video"
}

func (c videoClassifier) Priority() int {
	return 1
}

func (c videoClassifier) Classify(ctx context.Context, t model.Torrent) (classifier.Classification, error) {
	if hasVideo := t.HasFileType(model.FileTypeVideo); hasVideo.Valid && !hasVideo.Bool {
		return classifier.Classification{}, classifier.ErrNoMatch
	}
	if !t.Hint.IsNil() && !t.Hint.ContentType.IsVideo() {
		return classifier.Classification{}, classifier.ErrNoMatch
	}
	ct, title, year, attrs, err := ParseContent(t.Hint.NullContentType(), t.Name)
	if err != nil {
		return classifier.Classification{}, err
	}
	ref := t.Hint.ContentRef()
	if t.Hint.Title.Valid {
		title = t.Hint.Title.String
	}
	cl := classifier.Classification{
		ContentAttributes: attrs,
	}
	if content, err := c.resolveContent(ctx, ct, ref, title, year); err == nil {
		cl.Content = &content
	} else if !errors.Is(err, classifier.ErrNoMatch) {
		return classifier.Classification{}, err
	}
	cl.ApplyHint(t.Hint)
	if cl.Content != nil {
		cl.ContentType = model.NewNullContentType(cl.Content.Type)
		if cl.Content.OriginalLanguage.Valid {
			if len(cl.Languages) == 0 || cl.LanguageMulti {
				if cl.Languages == nil {
					cl.Languages = make(model.Languages)
				}
				cl.Languages[cl.Content.OriginalLanguage.Language] = struct{}{}
			}
		}
	}
	if !cl.ContentType.Valid {
		return classifier.Classification{}, classifier.ErrNoMatch
	}
	return cl, nil
}

func (c videoClassifier) resolveContent(
	ctx context.Context,
	ct model.ContentType,
	ref model.Maybe[model.ContentRef],
	title string,
	year model.Year,
) (model.Content, error) {
	if ct == model.ContentTypeMovie || ct == model.ContentTypeXxx {
		if ref.Valid {
			return c.tmdbClient.GetMovieByExternalId(ctx, ref.Val.Source, ref.Val.ID)
		}
		return c.tmdbClient.SearchMovie(ctx, tmdb.SearchMovieParams{
			Title:                title,
			Year:                 year,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		})
	}
	if ct == model.ContentTypeTvShow {
		if ref.Valid {
			return c.tmdbClient.GetTvShowByExternalId(ctx, ref.Val.Source, ref.Val.ID)
		}
		return c.tmdbClient.SearchTvShow(ctx, tmdb.SearchTvShowParams{
			Name:                 title,
			FirstAirDateYear:     year,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		})
	}
	return model.Content{}, classifier.ErrNoMatch
}
