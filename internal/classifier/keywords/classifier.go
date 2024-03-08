package keywords

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"regexp"
)

type keywordsClassifier struct {
	contentType       model.ContentType
	priority          int
	regex             *regexp.Regexp
	requiredFileTypes []model.FileType
}

func (c keywordsClassifier) Key() string {
	return "keywords_" + c.contentType.String()
}

func (c keywordsClassifier) Priority() int {
	return c.priority
}

func (c keywordsClassifier) Classify(_ context.Context, t model.Torrent) (classifier.Classification, error) {
	if !t.Hint.IsNil() || !c.regex.MatchString(t.Name) {
		return classifier.Classification{}, classifier.ErrNoMatch
	}
	if len(c.requiredFileTypes) > 0 {
		hasRequiredFileTypes := t.HasFileType(c.requiredFileTypes...)
		if hasRequiredFileTypes.Valid && !hasRequiredFileTypes.Bool {
			return classifier.Classification{}, classifier.ErrNoMatch
		}
	}
	cl := classifier.Classification{
		ContentType: model.NullContentType{
			Valid:       true,
			ContentType: c.contentType,
		},
	}
	hasVideo := t.HasFileType(model.FileTypeVideo)
	if hasVideo.Valid && hasVideo.Bool {
		cl.InferVideoAttributes(t.Name)
	}
	return cl, nil
}
