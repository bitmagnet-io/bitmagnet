package torznab

import "encoding/xml"

// Error represents an error that should be encoded to an XML response according to the Torznab specification.
type Error struct {
	XMLName xml.Name `xml:"error"`

	// Code the Error code.
	//
	// From https://torznab.github.io/spec-1.3-draft/external/newznab/api.html#newznab-error-codes:
	//
	// 100 Incorrect user credentials
	// 101 Account suspended
	// 102 Insufficient privileges/not authorized
	// 103 Registration denied
	// 104 Registrations are closed
	// 105 Invalid registration (Email Address Taken)
	// 106 Invalid registration (Email Address Bad Format)
	// 107 Registration Failed (Data error)
	// 200 Missing parameter
	// 201 Incorrect parameter
	// 202 No such function. (Function not defined in this specification).
	// 203 Function not available. (Optional function is not implemented).
	// 300 No such item.
	// 300 Item already exists.
	// 900 Unknown error
	// 910 API Disabled
	Code        int    `xml:"error,attr"`
	Description string `xml:"description,attr"`
}

func (e Error) Error() string {
	return e.Description
}

func (e Error) XML() ([]byte, error) {
	return objToXML(e)
}
