package search

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"gopkg.in/yaml.v3"
)

type Criteria interface {
	criteria()
}

type jsonSpec = json_spec.ParseContext

type compilerContext struct {
	jsonSpec
	criteria   []criteriaDefinition
	resultType *ResultType
}

func (c compilerContext) child(pathPart string, source any) compilerContext {
	c.jsonSpec = c.Child(pathPart, source)

	return c
}

type criteriaDefinition interface {
	json_spec.HasJSONSchema
	name() string
	compile(ctx compilerContext) (Criteria, error)
	resultTypes() map[ResultType]struct{}
}

var criteriaDefinitions = []criteriaDefinition{
	definitionAnd{},
	definitionOr{},
	definitionNot{},
	definitionContentType{},
	definitionContentRef{},
	definitionGenre{},
	definitionInfoHash{},
	definitionLanguage{},
	definitionQueryString{},
	definitionTag{},
}

func ParseCriteria(value []byte) (Criteria, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(value, &node); err != nil {
		return nil, err
	}

	return CompileCriteria(json_schema.JSONValue(node))
}

func CompileCriteria(value json_schema.JSONValue) (Criteria, error) {
	return compileCriteria(value, nil)
}

func CompileResultTypeCriteria(resultType ResultType, value json_schema.JSONValue) (Criteria, error) {
	return compileCriteria(value, &resultType)
}

func compileCriteria(value json_schema.JSONValue, resultType *ResultType) (Criteria, error) {
	ctx := compilerContext{
		jsonSpec: json_spec.ParseContext{
			KeyMatcher: json_spec.KeyMatcherLowerCamel,
			Source:     value.Raw(),
		},
		criteria:   criteriaDefinitions,
		resultType: resultType,
	}

	criteria, err := compileCriteriaCtx(ctx)
	if err != nil {
		return nil, err
	}

	if len(criteria) == 1 {
		return criteria[0], nil
	}

	return And(criteria), nil
}

func compileCriteriaCtx(ctx compilerContext) ([]Criteria, error) {
	var (
		rawCriteria []any
		criteria    []Criteria
		errs        []error
	)

	isArray := false

	if s, ok := ctx.Source.([]any); ok {
		rawCriteria = s
		isArray = true
	} else if ctx.Source != nil {
		rawCriteria = []any{ctx.Source}
	}

outer:
	for i, raw := range rawCriteria {
		actionCtx := ctx
		if isArray {
			actionCtx = ctx.child(json_spec.NumericPathPart(i), raw)
		}

		for _, def := range ctx.criteria {
			if ctx.resultType != nil {
				if _, ok := def.resultTypes()[*ctx.resultType]; !ok {
					continue
				}
			}

			c, err := def.compile(actionCtx.child(def.name(), raw))
			if err == nil {
				criteria = append(criteria, c)
				continue outer
			}

			if json_spec.AsFatalCompilerError(err) != nil {
				return nil, err
			}
		}

		errs = append(errs, fmt.Errorf("no criteria matched: %v", ctx.Source))
	}

	if len(errs) > 0 {
		return nil, ctx.Fatal(errors.Join(errs...))
	}

	return criteria, nil
}

var criteriaSpec = json_spec.MustSucceed[[]any]{
	Typed: json_spec.List[any]{
		ItemSpec: json_spec.Generic[any]{
			Schema: json_schema.MustNew(
				json_schema.RefDefinition("criteria"),
			),
		},
	},
}

var stringSpec = json_spec.Generic[string]{
	Schema: json_schema.MustNew(
		json_schema.Typed(json_schema.TypeString),
	),
}
