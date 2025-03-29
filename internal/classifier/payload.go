package classifier

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type HasJSONSchema interface {
	JSONSchema() JSONSchema
}

type TypedPayload[T any] interface {
	HasJSONSchema
	Unmarshal(ctx compilerContext) (T, error)
}

type PayloadTransformerFunc[From any, To any] func(From, compilerContext) (To, error)

type payloadTransformer[From any, To any] struct {
	spec      TypedPayload[From]
	transform PayloadTransformerFunc[From, To]
}

func (s payloadTransformer[From, To]) JSONSchema() JSONSchema {
	return s.spec.JSONSchema()
}

func (s payloadTransformer[From, To]) Unmarshal(ctx compilerContext) (to To, _ error) {
	from, err := s.spec.Unmarshal(ctx)
	if err != nil {
		return to, err
	}

	return s.transform(from, ctx)
}

type payloadUnion[T any] struct {
	oneOf []TypedPayload[T]
}

func (s payloadUnion[T]) JSONSchema() JSONSchema {
	schemas := make([]any, len(s.oneOf))
	for i, spec := range s.oneOf {
		schemas[i] = spec.JSONSchema()
	}

	return map[string]any{
		"oneOf": schemas,
	}
}

func (s payloadUnion[T]) Unmarshal(ctx compilerContext) (to T, _ error) {
	//nolint:prealloc
	var errs []error

	for _, def := range s.oneOf {
		result, err := def.Unmarshal(ctx)
		if err == nil {
			return result, nil
		}

		errs = append(errs, err)
	}

	errs = append(errs, errors.New("no definition matched"))

	return to, errors.Join(errs...)
}

type payloadGeneric[T any] struct {
	jsonSchema map[string]any
}

func (s payloadGeneric[T]) JSONSchema() JSONSchema {
	return s.jsonSchema
}

func (payloadGeneric[T]) Unmarshal(ctx compilerContext) (to T, err error) {
	to, ok := ctx.source.(T)
	if !ok {
		err = ctx.error(errors.New("not ok"))
	}

	return to, err
}

type payloadStruct[T any] struct {
	jsonSchema map[string]any
}

func (s payloadStruct[T]) JSONSchema() JSONSchema {
	return s.jsonSchema
}

func (payloadStruct[T]) Unmarshal(ctx compilerContext) (to T, err error) {
	return decode[T](ctx)
}

type payloadLiteral[T comparable] struct {
	literal     T
	description string
}

func (s payloadLiteral[T]) JSONSchema() JSONSchema {
	schema := map[string]any{
		"const": s.literal,
	}
	if s.description != "" {
		schema["description"] = s.description
	}

	return schema
}

func (s payloadLiteral[T]) Unmarshal(ctx compilerContext) (to T, _ error) {
	typedPayload, err := decode[T](ctx)
	if err != nil {
		return to, err
	}

	if typedPayload != s.literal {
		return to, errors.New("value mismatch")
	}

	return typedPayload, nil
}

type payloadList[T any] struct {
	itemSpec    TypedPayload[T]
	description string
}

func (s payloadList[T]) JSONSchema() JSONSchema {
	schema := map[string]any{
		"type":  "array",
		"items": s.itemSpec.JSONSchema(),
	}
	if s.description != "" {
		schema["description"] = s.description
	}

	return schema
}

func (s payloadList[T]) Unmarshal(ctx compilerContext) (to []T, _ error) {
	if ctx.source == nil {
		return nil, nil
	}

	rawList, ok := ctx.source.([]any)
	if !ok {
		rawList = []any{ctx.source}
	}

	to = make([]T, len(rawList))

	for i, rawItem := range rawList {
		item, err := s.itemSpec.Unmarshal(ctx.child(numericPathPart(i), rawItem))
		if err != nil {
			return to, err
		}

		to[i] = item
	}

	return to, nil
}

type payloadSingleKeyValue[T any] struct {
	key         string
	valueSpec   TypedPayload[T]
	description string
}

func (s payloadSingleKeyValue[T]) JSONSchema() JSONSchema {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			s.key: s.valueSpec.JSONSchema(),
		},
		"required":             []string{s.key},
		"additionalProperties": false,
	}
	if s.description != "" {
		schema["description"] = s.description
	}

	return schema
}

func (s payloadSingleKeyValue[T]) Unmarshal(ctx compilerContext) (to T, _ error) {
	rawMap, err := decode[map[string]any](ctx)
	if err != nil {
		return to, err
	}

	if len(rawMap) != 1 {
		return to, ctx.error(errors.New("expected a single key"))
	}

	rawValue, ok := rawMap[s.key]
	if !ok {
		return to, ctx.error(fmt.Errorf("missing expected key: '%s' %+v", s.key, rawMap))
	}

	value, err := s.valueSpec.Unmarshal(ctx.child(s.key, rawValue))
	if err != nil {
		return to, err
	}

	return value, nil
}

type payloadEnum[T string] struct {
	values []T
}

func (s payloadEnum[T]) JSONSchema() JSONSchema {
	return map[string]any{
		"type": "string",
		"enum": s.values,
	}
}

func (s payloadEnum[T]) Unmarshal(ctx compilerContext) (to T, _ error) {
	value, err := decode[T](ctx)
	if err != nil {
		return to, ctx.error(err)
	}

	for _, validValue := range s.values {
		if value == validValue {
			return value, nil
		}
	}

	return to, ctx.error(fmt.Errorf("value not in enum: '%s'", value))
}

type payloadMustSucceed[T any] struct {
	payload TypedPayload[T]
}

func (p payloadMustSucceed[T]) Unmarshal(ctx compilerContext) (t T, _ error) {
	result, err := p.payload.Unmarshal(ctx)
	if err != nil {
		return t, ctx.fatal(err)
	}

	return result, nil
}

func (p payloadMustSucceed[T]) JSONSchema() JSONSchema {
	return p.payload.JSONSchema()
}

var contentTypePayloadSpec = payloadTransformer[string, model.NullContentType]{
	spec: payloadEnum[string]{append(model.ContentTypeNames(), "unknown")},
	transform: func(str string, _ compilerContext) (model.NullContentType, error) {
		if str == "unknown" {
			return model.NullContentType{}, nil
		}
		contentType, err := model.ParseContentType(str)
		if err != nil {
			return model.NullContentType{}, err
		}
		return model.NullContentType{ContentType: contentType, Valid: true}, nil
	},
}
