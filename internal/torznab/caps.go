package torznab

import "encoding/xml"

type Caps struct {
	XMLName    xml.Name       `xml:"caps"`
	Server     CapsServer     `xml:"server"`
	Limits     CapsLimits     `xml:"limits"`
	Searching  CapsSearching  `xml:"searching"`
	Categories CapsCategories `xml:"categories"`
	Tags       string         `xml:"tags"`
}

func (c Caps) XML() ([]byte, error) {
	return objToXML(c)
}

type CapsServer struct {
	Title string `xml:"title,attr"`
}

type CapsLimits struct {
	Max     uint `xml:"max,attr,omitempty"`
	Default uint `xml:"default,attr,omitempty"`
}

type CapsSearch struct {
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr,omitempty"`
}

type CapsSearching struct {
	Search      CapsSearch `xml:"search,omitempty"`
	TvSearch    CapsSearch `xml:"tv-search,omitempty"`
	MovieSearch CapsSearch `xml:"movie-search,omitempty"`
	MusicSearch CapsSearch `xml:"music-search,omitempty"`
	AudioSearch CapsSearch `xml:"audio-search,omitempty"`
	BookSearch  CapsSearch `xml:"book-search,omitempty"`
}

type CapsCategories struct {
	Categories []Category `xml:"category"`
}
