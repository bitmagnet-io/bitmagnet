package reject

import (
	"github.com/pemistahl/lingua-go"
)

func IsLanguageNotAllowed(title string) bool {
	languages := []lingua.Language{
		lingua.English,
		lingua.French,
		lingua.German,
		lingua.Spanish,
		lingua.Italian,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	_, exists := detector.DetectLanguageOf(title)
	return !exists
}
