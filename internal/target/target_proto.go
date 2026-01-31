package target

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_forms"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/transformer/transform_model"
	"github.com/xeipuuv/gojsonschema"
)

type targetProto struct {
	ref            ref.Ref
	name           string
	impl           api.TorrentTarget
	dataSchemaOnce sync.Once
	dataSchema     *dataSchema
	dataSchemaErr  error
}

func NewTargetProto(
	ref ref.Ref,
	name string,
	impl api.TorrentTarget,
) TorrentContentTarget {
	return &targetProto{
		ref:  ref,
		name: name,
		impl: impl,
	}
}

func (t *targetProto) Ref() ref.Ref {
	return t.ref
}

func (t *targetProto) Name() string {
	return t.name
}

func (t *targetProto) DataSchama(ctx context.Context) (*json_schema.JSONSchema, error) {
	schema, err := t.dataSchama(ctx)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		//nolint:nilnil
		return nil, nil
	}

	return schema.jsonSchema, nil
}

func (t *targetProto) dataSchama(ctx context.Context) (*dataSchema, error) {
	t.dataSchemaOnce.Do(func() {
		res, err := t.impl.DataSchema(ctx, &api.Empty{})
		if err != nil {
			t.dataSchemaErr = err
			return
		}

		if res == nil || len(res.GetData()) == 0 {
			return
		}

		var js json_schema.JSONSchema
		if err := js.UnmarshalJSON(res.GetData()); err != nil {
			t.dataSchemaErr = err
			return
		}

		var gjsSchema *gojsonschema.Schema

		gjsSchema, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(res.GetData()))
		if err != nil {
			t.dataSchemaErr = err
			return
		}

		t.dataSchema = &dataSchema{
			jsonSchema: &js,
			gjsSchema:  gjsSchema,
		}
	})

	if t.dataSchemaErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrDataSchema, t.dataSchemaErr)
	}

	return t.dataSchema, nil
}

// todo: Cache
func (t *targetProto) UISchema(ctx context.Context, acceptLanguage []string) (*json_forms.UISchema, error) {
	res, err := t.impl.UISchema(ctx, &api.LocalizeParams{
		AcceptLanguage: acceptLanguage,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUISchema, err)
	}

	if res == nil || len(res.GetData()) == 0 {
		//nolint:nilnil
		return nil, nil
	}

	uiSchema, err := json_forms.UnmarshalUISchema(res.GetData())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUISchema, err)
	}

	return &uiSchema, nil
}

func (t *targetProto) Send(
	ctx context.Context,
	torrents []model.TorrentContent,
	data json_schema.JSONValue,
) (*json_schema.JSONValue, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMarshalData, err)
	}

	schema, err := t.dataSchama(ctx)
	if err != nil {
		return nil, err
	}

	if schema != nil {
		if err := schema.validate(jsonBytes); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrValidation, err)
		}
	}

	req := &api.SendTorrentsParams{
		Torrents: slice.Map(torrents, transform_model.TorrentContentToProto),
		Data:     jsonBytes,
	}

	res, err := t.impl.Send(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPlugin, err)
	}

	var returnData *json_schema.JSONValue

	if len(res.GetData()) > 0 {
		var raw any

		err = json.Unmarshal(res.GetData(), &raw)
		if err == nil {
			val, err := json_schema.NewValue(raw)
			if err == nil {
				returnData = &val
			}
		}
	}

	return returnData, nil
}

type dataSchema struct {
	jsonSchema *json_schema.JSONSchema
	gjsSchema  *gojsonschema.Schema
}

func (s *dataSchema) validate(jsonBytes []byte) error {
	if result, err := s.gjsSchema.Validate(gojsonschema.NewBytesLoader(jsonBytes)); err != nil {
		return err
	} else if !result.Valid() {
		return errors.New(strings.Join(slice.Map(result.Errors(), func(e gojsonschema.ResultError) string {
			return e.String()
		}), "; "))
	}

	return nil
}
