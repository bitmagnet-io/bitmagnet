package model

import "gorm.io/gorm"

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
	Url string
}

func (c Content) ExternalLinks() []ExternalLink {
	links := make([]ExternalLink, 0)
	if link := getExternalLinkUrl(c.Type, c.Source, c.ID); link.Valid {
		links = append(links, ExternalLink{
			MetadataSource: c.MetadataSource,
			Url:            link.String,
		})
	}
	for _, attr := range c.Attributes {
		if attr.Key == "id" {
			if link := getExternalLinkUrl(c.Type, attr.Source, attr.Value); link.Valid {
				links = append(links, ExternalLink{
					MetadataSource: attr.MetadataSource,
					ID:             attr.Value,
					Url:            link.String,
				})
			}
		}
	}
	return links
}

func getExternalLinkUrl(contentType ContentType, source, id string) NullString {
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

func (c *Content) BeforeSave(tx *gorm.DB) error {
	c.SearchString = c.Title + " " + c.ReleaseDate.YearString()
	if c.OriginalTitle.Valid && c.Title != c.OriginalTitle.String {
		c.SearchString += " " + c.OriginalTitle.String
	}
	for _, collection := range c.Collections {
		c.SearchString += " " + collection.Name
	}
	return nil
}
