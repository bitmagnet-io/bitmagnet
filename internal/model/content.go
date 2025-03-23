package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
)

type ContentRef struct {
	Type   ContentType
	Source string
	ID     string
}

func (c Content) Ref() ContentRef {
	return ContentRef{
		Type:   c.Type,
		Source: c.Source,
		ID:     c.ID,
	}
}

func (c Content) Identifier(source string) (string, bool) {
	if c.Source == source {
		return c.ID, true
	}

	for _, attr := range c.Attributes {
		if attr.Key == "id" && attr.Source == source {
			return attr.Value, true
		}
	}

	return "", false
}

type ExternalLink struct {
	MetadataSource
	ID  string
	URL string
}

func (c Content) ExternalLinks() []ExternalLink {
	links := make([]ExternalLink, 0)
	if link := getExternalLinkURL(c.Type, c.Source, c.ID); link.Valid {
		links = append(links, ExternalLink{
			MetadataSource: c.MetadataSource,
			URL:            link.String,
		})
	}

	for _, attr := range c.Attributes {
		if attr.Key == "id" {
			if link := getExternalLinkURL(c.Type, attr.Source, attr.Value); link.Valid {
				links = append(links, ExternalLink{
					MetadataSource: attr.MetadataSource,
					ID:             attr.Value,
					URL:            link.String,
				})
			}
		}
	}

	return links
}

func getExternalLinkURL(contentType ContentType, source, id string) NullString {
	switch source {
	case "imdb":
		return NewNullString("https://www.imdb.com/title/" + id)
	case "tmdb":
		switch contentType {
		case ContentTypeTvShow:
			return NewNullString("https://www.themoviedb.org/tv/" + id)
		default:
			return NewNullString("https://www.themoviedb.org/movie/" + id)
		}
	case "tvdb":
		return NewNullString("https://www.thetvdb.com/dereferrer/series/" + id)
	}

	return NullString{}
}

func (c *Content) UpdateTsv() {
	tsv := fts.Tsvector{}
	tsv.AddText(c.Title, fts.TsvectorWeightA)

	if c.OriginalTitle.Valid && c.Title != c.OriginalTitle.String {
		tsv.AddText(c.OriginalTitle.String, fts.TsvectorWeightA)
	}

	if !c.ReleaseYear.IsNil() {
		tsv.AddText(c.ReleaseYear.String(), fts.TsvectorWeightB)
	}

	for _, c := range c.Collections {
		if c.Type == "genre" {
			tsv.AddText(c.Name, fts.TsvectorWeightD)
		}
	}

	for _, a := range c.Attributes {
		if a.Key == "id" {
			tsv.AddText(a.Value, fts.TsvectorWeightD)
		}
	}

	c.Tsv = tsv
}
