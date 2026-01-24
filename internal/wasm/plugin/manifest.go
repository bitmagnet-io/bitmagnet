package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/xeipuuv/gojsonschema"
)

type Manifest struct {
	Name         string                            `json:"name"`
	Description  string                            `json:"description"`
	Version      string                            `json:"version"`
	Config       map[string]json_schema.JSONSchema `json:"config"`
	Concurrency  int                               `json:"concurrency"`
	Capabilities Capabilities                      `json:"capabilities"`
	Permissions  Permissions                       `json:"permissions"`
}

func ParseManifest(data []byte) (Manifest, error) {
	loader := gojsonschema.NewBytesLoader(data)

	result, err := manifestGoJSONSchema.Validate(loader)
	if err != nil {
		return Manifest{}, err
	}

	if !result.Valid() {
		var errs []string
		for _, desc := range result.Errors() {
			errs = append(errs, desc.String())
		}

		return Manifest{}, fmt.Errorf("invalid manifest: %s", strings.Join(errs, "; "))
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}
