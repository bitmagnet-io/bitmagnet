package torznab

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

type SearchResult struct {
	XMLName    xml.Name            `xml:"rss"`
	RssVersion rssVersion          `xml:"version,attr"`
	AtomNs     customNs            `xml:"xmlns:atom,attr"`
	TorznabNs  customNs            `xml:"xmlns:torznab,attr"`
	Channel    SearchResultChannel `xml:"channel"`
}

func (r SearchResult) Xml() ([]byte, error) {
	return objToXml(r)
}

type customNs struct{}

func (r customNs) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	switch name.Local {
	case "xmlns:atom":
		return xml.Attr{
			Name:  name,
			Value: "http://www.w3.org/2005/Atom",
		}, nil
	case "xmlns:torznab":
		return xml.Attr{
			Name:  name,
			Value: "http://torznab.com/schemas/2015/feed",
		}, nil
	default:
		return xml.Attr{}, errors.New("unknown namespace")
	}
}

type rssVersion string

func (r rssVersion) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: "2.0",
	}, nil
}

type RssDate time.Time

const RssDateDefaultFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

func (r RssDate) String() string {
	return time.Time(r).Format(RssDateDefaultFormat)
}

var rssDateFormats = []string{
	RssDateDefaultFormat,
	// if parsing is needed in future we might need these:
	//time.RFC850,
	//time.RFC822,
	//time.RFC822Z,
	//time.RFC1123,
	//time.RFC1123Z,
	//"02 Jan 2006 15:04:05 MST",
	//"02 Jan 2006 15:04:05 -0700",
}

func (r *RssDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	for _, format := range rssDateFormats {
		parsed, err := time.Parse(format, v)
		if err == nil {
			*r = RssDate(parsed)
			return nil
		}
	}
	return fmt.Errorf("cannot parse %q as RssDate", v)
}

func (r RssDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(r.String(), start)
}

type SearchResultChannel struct {
	Title         string               `xml:"title,omitempty"`
	Link          string               `xml:"link,omitempty"`
	Description   string               `xml:"description,omitempty"`
	Language      string               `xml:"language,omitempty"`
	PubDate       RssDate              `xml:"pubDate,omitempty"`
	LastBuildDate RssDate              `xml:"lastBuildDate,omitempty"`
	Docs          string               `xml:"docs,omitempty"`
	Generator     string               `xml:"generator,omitempty"`
	Response      SearchResultResponse `xml:"http://www.newznab.com/DTD/2010/feeds/attributes/ response"`
	Items         []SearchResultItem   `xml:"item"`
}

type SearchResultResponse struct {
	Offset uint `xml:"offset,attr,omitempty"`
	Total  uint `xml:"total,attr,omitempty"`
}

type SearchResultItem struct {
	Title        string                        `xml:"title"`
	GUID         string                        `xml:"guid,omitempty"`
	PubDate      RssDate                       `xml:"pubDate,omitempty"`
	Category     string                        `xml:"category,omitempty"`
	Link         string                        `xml:"link,omitempty"`
	Size         uint                          `xml:"size"`
	Description  string                        `xml:"description,omitempty"`
	Comments     string                        `xml:"comments,omitempty"`
	Enclosure    SearchResultItemEnclosure     `xml:"enclosure"`
	TorznabAttrs []SearchResultItemTorznabAttr `xml:"torznab:attr"`
}

type SearchResultItemEnclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type SearchResultItemTorznabAttr struct {
	AttrName  string `xml:"name,attr"`
	AttrValue string `xml:"value,attr"`
}
