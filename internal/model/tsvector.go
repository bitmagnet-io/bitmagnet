package model

// TsvectorLabel represents the label used for weightings of search parts
// See https://www.postgresql.org/docs/current/textsearch-controls.html#TEXTSEARCH-PARSING-DOCUMENTS
// ENUM(A, B, C, D)
type TsvectorLabel string

type TsvectorPositionLabel struct {
	Position int
	Label    TsvectorLabel
}

type Tsvector map[string]map[TsvectorPositionLabel]struct{}

func NewTsvParts() Tsvector {
	return Tsvector{}
}

func (t Tsvector) Add(label TsvectorLabel, value string) {
	nextPos := 1
	for _, pls := range t {
		for pl := range pls {
			if pl.Label == label && pl.Position >= nextPos {
				nextPos = pl.Position + 1
			}
		}
	}
	for _, word := range fts.TokenizeFlat(value) {
		if _, ok := t[word]; !ok {
			t[word] = make(map[TsvectorPositionLabel]struct{})
		}
		t[word][TsvectorPositionLabel{
			Position: nextPos,
			Label:    label,
		}] = struct{}{}
		nextPos++
	}
}
