package json_spec

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type HasJSONSchema interface {
	JSONSchema() json_schema.JSONSchema
}

type Typed[T any] interface {
	HasJSONSchema
	Parse(ctx ParseContext) (T, error)
}

type TransformerFunc[From any, To any] func(From, ParseContext) (To, error)

type Transformer[From any, To any] struct {
	Typed[From]
	Transform TransformerFunc[From, To]
}

func (s Transformer[From, To]) Parse(ctx ParseContext) (to To, _ error) {
	from, err := s.Typed.Parse(ctx)
	if err != nil {
		return to, err
	}

	return s.Transform(from, ctx)
}

type Union[T any] struct {
	OneOf []Typed[T]
}

func (s Union[T]) JSONSchema() json_schema.JSONSchema {
	schemas := make([]json_schema.JSONSchema, len(s.OneOf))
	for i, spec := range s.OneOf {
		schemas[i] = spec.JSONSchema()
	}

	return json_schema.MustNew(json_schema.OneOf(schemas...))
}

func (s Union[T]) Parse(ctx ParseContext) (to T, _ error) {
	var errs []error

	for _, def := range s.OneOf {
		result, err := def.Parse(ctx)
		if err == nil {
			return result, nil
		}

		errs = append(errs, err)
	}

	errs = append(errs, errors.New("no definition matched"))

	return to, errors.Join(errs...)
}

type Generic[T any] struct {
	Schema json_schema.JSONSchema
	Parser func(ctx ParseContext) (T, error)
}

func (s Generic[T]) JSONSchema() json_schema.JSONSchema {
	return s.Schema
}

func (s Generic[T]) Parse(ctx ParseContext) (to T, err error) {
	if s.Parser != nil {
		return s.Parser(ctx)
	}

	to, ok := ctx.Source.(T)
	if !ok {
		err = ctx.Error(errors.New("not ok"))
	}

	return to, err
}

type Struct[T any] struct {
	Schema json_schema.JSONSchema
}

func (s Struct[T]) JSONSchema() json_schema.JSONSchema {
	return s.Schema
}

func (Struct[T]) Parse(ctx ParseContext) (to T, err error) {
	return Decode[T](ctx)
}

type Literal[T comparable] struct {
	Literal     T
	Description string
}

func (s Literal[T]) JSONSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Const(json_schema.MustNewValue(s.Literal)),
		json_schema.DescriptionIfNonEmpty(s.Description),
	)
}

func (s Literal[T]) Parse(ctx ParseContext) (to T, _ error) {
	typedPayload, err := Decode[T](ctx)
	if err != nil {
		return to, err
	}

	if typedPayload != s.Literal {
		return to, errors.New("value mismatch")
	}

	return typedPayload, nil
}

type List[T any] struct {
	ItemSpec    Typed[T]
	Description string
}

func (s List[T]) JSONSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeArray),
		json_schema.Items(s.ItemSpec.JSONSchema()),
		json_schema.DescriptionIfNonEmpty(s.Description),
	)
}

func (s List[T]) Parse(ctx ParseContext) (to []T, _ error) {
	if ctx.Source == nil {
		return nil, nil
	}

	rawList, ok := ctx.Source.([]any)
	if !ok {
		rawList = []any{ctx.Source}
	}

	to = make([]T, len(rawList))

	for i, rawItem := range rawList {
		item, err := s.ItemSpec.Parse(ctx.Child(NumericPathPart(i), rawItem))
		if err != nil {
			return nil, err
		}

		to[i] = item
	}

	return to, nil
}

type SingleKeyValue[T any] struct {
	Key         string
	ValueSpec   Typed[T]
	Description string
}

func (s SingleKeyValue[T]) JSONSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
		json_schema.Properties(map[string]json_schema.JSONSchema{
			s.Key: s.ValueSpec.JSONSchema(),
		}),
		json_schema.Required(json_schema.RequiredFields{s.Key}),
		json_schema.AdditionalPropertiesFalse(),
		json_schema.DescriptionIfNonEmpty(s.Description),
	)
}

func (s SingleKeyValue[T]) Parse(ctx ParseContext) (to T, _ error) {
	rawMap, err := Decode[map[string]any](ctx)
	if err != nil {
		return to, err
	}

	if len(rawMap) != 1 {
		return to, ctx.Error(errors.New("expected a single key"))
	}

	rawValue, ok := rawMap[s.Key]
	if !ok {
		return to, ctx.Error(fmt.Errorf("missing expected key: '%s' %+v", s.Key, rawMap))
	}

	value, err := s.ValueSpec.Parse(ctx.Child(s.Key, rawValue))
	if err != nil {
		return to, err
	}

	return value, nil
}

type Enum[T ~string] struct {
	Values []T
}

func (s Enum[T]) JSONSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeString),
		json_schema.Enum(slice.Map(s.Values, func(v T) json_schema.JSONValue {
			return json_schema.MustNewValue(string(v))
		})...),
	)
}

func (s Enum[T]) Parse(ctx ParseContext) (to T, _ error) {
	value, err := Decode[T](ctx)
	if err != nil {
		return to, ctx.Error(err)
	}

	for _, validValue := range s.Values {
		if value == validValue {
			return value, nil
		}
	}

	return to, ctx.Error(fmt.Errorf("value not in enum: '%s'", value))
}

type MustSucceed[T any] struct {
	Typed[T]
}

func (p MustSucceed[T]) Parse(ctx ParseContext) (t T, _ error) {
	result, err := p.Typed.Parse(ctx)
	if err != nil {
		return t, ctx.Fatal(err)
	}

	return result, nil
}
