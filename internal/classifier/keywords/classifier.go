package keywords

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"regexp"
)

func NewKeywordsClassifier(contentType model.ContentType, words []string, priority int) classifier.SubClassifier {
	return keywordsClassifier{
		contentType: contentType,
		priority:    priority,
		regex:       regex.NewRegexFromNames(words...),
	}
}

type keywordsClassifier struct {
	contentType model.ContentType
	priority    int
	regex       *regexp.Regexp
}

func (c keywordsClassifier) Key() string {
	return "keywords_" + c.contentType.String()
}

func (c keywordsClassifier) Priority() int {
	return c.priority
}

func (c keywordsClassifier) Classify(_ context.Context, torrent model.Torrent) (classifier.Classification, error) {
	if !c.regex.MatchString(torrent.Name) {
		return classifier.Classification{}, classifier.ErrNoMatch
	}
	return classifier.Classification{
		ContentType: model.NullContentType{
			Valid:       true,
			ContentType: c.contentType,
		},
	}, nil
}
