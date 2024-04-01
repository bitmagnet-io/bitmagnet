package workflow

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type HasJsonSchema interface {
	JsonSchema() any
}

type TypedPayload[T any] interface {
	HasJsonSchema
	Unmarshal(ctx compilerContext) (T, error)
}

type PayloadTransformerFunc[From any, To any] func(From, compilerContext) (To, error)

type payloadTransformer[From any, To any] struct {
	spec      TypedPayload[From]
	transform PayloadTransformerFunc[From, To]
}

func (s payloadTransformer[From, To]) JsonSchema() any {
	return s.spec.JsonSchema()
}

func (s payloadTransformer[From, To]) Unmarshal(ctx compilerContext) (to To, _ error) {
	from, err := s.spec.Unmarshal(ctx)
	if err != nil {
		return to, err
	}
	return s.transform(from, ctx)
}

func payloadIdentityTransformer[From any, To any](value To) PayloadTransformerFunc[From, To] {
	return func(From, compilerContext) (To, error) {
		return value, nil
	}
}

type payloadUnion[T any] struct {
	oneOf []TypedPayload[T]
}

func (s payloadUnion[T]) JsonSchema() any {
	schemas := make([]any, len(s.oneOf))
	for i, spec := range s.oneOf {
		schemas[i] = spec.JsonSchema()
	}
	return map[string]any{
		"oneOf": schemas,
	}
}

func (s payloadUnion[T]) Unmarshal(ctx compilerContext) (to T, _ error) {
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
	jsonSchema any
}

func (s payloadGeneric[T]) JsonSchema() any {
	return s.jsonSchema
}

func (s payloadGeneric[T]) Unmarshal(ctx compilerContext) (to T, err error) {
	to, ok := ctx.source.(T)
	if !ok {
		err = ctx.error(errors.New("not ok"))
	}
	return to, err
}

type payloadStruct[T any] struct {
	jsonSchema any
}

func (s payloadStruct[T]) JsonSchema() any {
	return s.jsonSchema
}

func (s payloadStruct[T]) Unmarshal(ctx compilerContext) (to T, err error) {
	return decode[T](ctx)
}

type payloadLiteral[T comparable] struct {
	literal T
}

func (s payloadLiteral[T]) JsonSchema() any {
	return map[string]any{
		"const":    s.literal,
		"nullable": false,
	}
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
	itemSpec TypedPayload[T]
}

func (s payloadList[T]) JsonSchema() any {
	return map[string]any{
		"type":  "array",
		"items": s.itemSpec.JsonSchema(),
	}
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
	key       string
	valueSpec TypedPayload[T]
}

func (s payloadSingleKeyValue[T]) JsonSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			s.key: s.valueSpec.JsonSchema(),
		},
		"required":             []string{s.key},
		"additionalProperties": false,
		"nullable":             false,
	}
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

func (s payloadEnum[T]) JsonSchema() any {
	return map[string]any{
		"type":     "string",
		"enum":     s.values,
		"nullable": false,
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

func (p payloadMustSucceed[T]) JsonSchema() any {
	return p.payload.JsonSchema()
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

type payloadKeyValue[T any] struct {
	valueSpec TypedPayload[T]
}

func (s payloadKeyValue[T]) JsonSchema() any {
	return map[string]any{
		"type": "object",
		"additionalProperties": map[string]any{
			"type": s.valueSpec.JsonSchema(),
		},
		"nullable": false,
	}
}

func (s payloadKeyValue[T]) Unmarshal(ctx compilerContext) (to map[string]T, _ error) {
	rawMap, err := decode[map[string]any](ctx)
	if err != nil {
		return to, err
	}
	kvs := make(map[string]T, len(rawMap))
	for key, rawValue := range rawMap {
		value, err := s.valueSpec.Unmarshal(ctx.child(key, rawValue))
		if err != nil {
			return to, err
		}
		kvs[key] = value
	}
	return kvs, nil
}
