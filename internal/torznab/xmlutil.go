package torznab

import (
	"encoding/xml"
)

type XMLer interface {
	XML() ([]byte, error)
}

func objToXML(obj any) ([]byte, error) {
	body, err := xml.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(body)), nil
}
