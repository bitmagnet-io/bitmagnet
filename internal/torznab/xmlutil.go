package torznab

import (
	"encoding/xml"
)

type Xmler interface {
	Xml() ([]byte, error)
}

func objToXml(obj any) ([]byte, error) {
	body, err := xml.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(body)), nil
}
