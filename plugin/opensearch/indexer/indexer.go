//go:build wasip1

package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/plugin-opensearch/shared"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func New(
	client *opensearchapi.Client,
	indexPrefix string,
) api.Indexer {
	return &indexer{
		client:           client,
		indexAliasPrefix: indexPrefix,
		indexPrefix:      fmt.Sprintf("%sv%d-", indexPrefix, templateVersion),
	}
}

type indexer struct {
	initialized      atomic.Bool
	client           *opensearchapi.Client
	indexPrefix      string
	indexAliasPrefix string
}

type indexID struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}

type actionUpdate struct {
	Update indexID `json:"update"`
}

func createLines(values ...any) ([]byte, error) {
	var data []byte

	for _, value := range values {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		data = append(data, jsonBytes...)
		data = append(data, '\n')
	}

	return data, nil
}

func (i *indexer) Index(ctx context.Context, payload *api.IndexPayload) (*api.Empty, error) {
	if i.initialized.CompareAndSwap(false, true) {
		templates, err := i.templates()
		if err != nil {
			return nil, err
		}

		for _, template := range templates {
			res, err := i.client.IndexTemplate.Create(ctx, template)
			if err == nil {
				err = responseError(res.Inspect().Response)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	var rawLines []any

	content := payload.GetContent()

	if len(content) > 0 {
		for _, c := range content {
			rawLines = append(rawLines,
				actionUpdate{Update: indexID{
					Index: i.indexPrefix + string(shared.RecordContent),
					ID:    c.GetRef().String(),
				}},
				map[string]any{
					"doc":           c,
					"doc_as_upsert": true,
				},
			)
		}
	}

	torrentContent := payload.GetTorrentContent()

	if len(torrentContent) > 0 {
		for _, tc := range torrentContent {
			action := actionUpdate{Update: indexID{
				Index: i.indexPrefix + string(shared.RecordTorrent),
				ID:    tc.GetId(),
			}}
			rawLines = append(rawLines,
				action,
				map[string]any{
					"doc":           tc,
					"doc_as_upsert": true,
				},
			)
		}
	}

	data, err := createLines(rawLines...)
	if err != nil {
		return nil, err
	}

	res, err := i.client.Bulk(ctx, opensearchapi.BulkReq{
		Body: bytes.NewReader(data),
	})

	if err == nil {
		err = responseError(res.Inspect().Response)
	}

	if err == nil {
		var errs []error
		for _, item := range res.Items {
			for _, resp := range item {
				if resp.Error != nil {
					errs = append(errs, fmt.Errorf("indexing error for ID %s: %s", resp.ID, resp.Error.Reason))
				}
			}
		}

		err = errors.Join(errs...)
	}

	return &api.Empty{}, err
}

func responseError(res *opensearch.Response) error {
	if res.StatusCode >= 300 {
		err := fmt.Errorf("unexpected status code: %d", res.StatusCode)
		body, _ := io.ReadAll(res.Body)
		if len(body) > 0 {
			err = fmt.Errorf("%w: %s", err, string(body))
		}
		return err
	}

	return nil
}

func newMappingType(typeName string) map[string]any {
	return map[string]any{
		"type": typeName,
	}
}
