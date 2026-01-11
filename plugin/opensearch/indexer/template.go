//go:build wasip1

package indexer

import (
	"bytes"
	"encoding/json"

	"github.com/bitmagnet-io/plugin-opensearch/shared"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const templateVersion = 1

func (i *indexer) templates() ([]opensearchapi.IndexTemplateCreateReq, error) {
	var templates []opensearchapi.IndexTemplateCreateReq

	if template, err := i.createTemplate(shared.RecordTorrent); err != nil {
		return nil, err
	} else {
		templates = append(templates, template)
	}

	if template, err := i.createTemplate(shared.RecordContent); err != nil {
		return nil, err
	} else {
		templates = append(templates, template)
	}

	return templates, nil
}

func (i *indexer) createTemplate(recordType shared.RecordType) (opensearchapi.IndexTemplateCreateReq, error) {
	body := map[string]any{
		"version":        templateVersion,
		"index_patterns": []string{i.indexPrefix + string(recordType)},
		"template": map[string]any{
			"aliases": map[string]any{
				i.indexAliasPrefix + string(recordType): map[string]any{},
			},
			"settings": map[string]any{
				"index": map[string]any{
					"number_of_shards":   1,
					"number_of_replicas": 1,
				},
			},
			"mappings": createMappings(recordType),
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return opensearchapi.IndexTemplateCreateReq{}, err
	}

	return opensearchapi.IndexTemplateCreateReq{
		IndexTemplate: i.indexPrefix + string(recordType),
		Body:          bytes.NewReader(data),
	}, nil
}

var (
	mappingKeyword = newMappingType("keyword")

	mappingLong = newMappingType("long")

	mappingInteger = newMappingType("integer")

	mappingFloat = newMappingType("float")

	mappingTimestamp = map[string]any{
		"type":   "date",
		"format": "epoch_millis",
	}

	mappingText = newMappingType("text")

	mappingTextWithKeyword = map[string]any{
		"type": "text",
		"fields": map[string]any{
			"keyword": mappingKeyword,
		},
	}

	mappingBoolean = newMappingType("boolean")
)

func createMappings(recordType shared.RecordType) map[string]any {
	switch recordType {
	case shared.RecordTorrent:
		return map[string]any{
			"dynamic": false,
			"properties": map[string]any{
				"id":          mappingKeyword,
				"infoHash":    mappingKeyword,
				"contentType": mappingKeyword,
				"contentRef": map[string]any{
					"properties": map[string]any{
						"type":   mappingKeyword,
						"source": mappingKeyword,
						"id":     mappingKeyword,
					},
				},
				"languages":       mappingKeyword,
				"episodes":        mappingKeyword,
				"videoResolution": mappingKeyword,
				"videoSource":     mappingKeyword,
				"size":            mappingLong,
				"filesCount":      mappingInteger,
				"seeders":         mappingInteger,
				"leechers":        mappingInteger,
				"publishedAt":     mappingTimestamp,
				"createdAt":       mappingTimestamp,
				"updatedAt":       mappingTimestamp,
				"torrent": map[string]any{
					"properties": map[string]any{
						"infoHash":    mappingKeyword,
						"name":        mappingTextWithKeyword,
						"size":        mappingLong,
						"filesStatus": mappingInteger,
						"extension":   mappingKeyword,
						"private":     mappingBoolean,
						"createdAt":   mappingTimestamp,
						"updatedAt":   mappingTimestamp,
						"files": map[string]any{
							"type": "nested",
							"properties": map[string]any{
								"index":     mappingInteger,
								"infoHash":  mappingKeyword,
								"path":      mappingTextWithKeyword,
								"extension": mappingKeyword,
								"fileType":  mappingKeyword,
								"size":      mappingLong,
								"createdAt": mappingTimestamp,
								"updatedAt": mappingTimestamp,
							},
						},
						"sources": map[string]any{
							"type": "nested",
							"properties": map[string]any{
								"infoHash":    mappingKeyword,
								"source":      mappingKeyword,
								"importId":    mappingKeyword,
								"seeders":     mappingInteger,
								"leechers":    mappingInteger,
								"publishedAt": mappingTimestamp,
								"createdAt":   mappingTimestamp,
								"updatedAt":   mappingTimestamp,
							},
						},
					},
				},
				"content": createMappings(shared.RecordContent),
			},
		}
	case shared.RecordContent:
		return map[string]any{
			"dynamic": false,
			"properties": map[string]any{
				"ref": map[string]any{
					"properties": map[string]any{
						"type":   mappingKeyword,
						"source": mappingKeyword,
						"id":     mappingKeyword,
					},
				},
				"title": map[string]any{
					"type": "text",
					"fields": map[string]any{
						"keyword": mappingKeyword,
					},
				},
				"releaseDate": map[string]any{
					"properties": map[string]any{
						"year":  mappingInteger,
						"month": mappingInteger,
						"day":   mappingInteger,
					},
				},
				"originalLanguage": mappingKeyword,
				"originalTitle":    mappingTextWithKeyword,
				"overview":         mappingText,
				"popularity":       mappingFloat,
				"voteAverage":      mappingFloat,
				"voteCount":        mappingInteger,
				"collections": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"type":      mappingKeyword,
						"source":    mappingKeyword,
						"id":        mappingKeyword,
						"name":      mappingKeyword,
						"createdAt": mappingTimestamp,
						"updatedAt": mappingTimestamp,
					},
				},
				"attributes": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"source":    mappingKeyword,
						"key":       mappingKeyword,
						"value":     mappingKeyword,
						"createdAt": mappingTimestamp,
						"updatedAt": mappingTimestamp,
					},
				},
				"createdAt": mappingTimestamp,
				"updatedAt": mappingTimestamp,
			},
		}
	default:
		return nil
	}
}
