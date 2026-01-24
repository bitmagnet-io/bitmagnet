package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const notName = "not"

type notCondition struct{}

func (notCondition) name() string {
	return notName
}

var notConditionSpec = json_spec.SingleKeyValue[any]{
	Key: notName,
	ValueSpec: json_spec.MustSucceed[any]{json_spec.Generic[any]{
		Schema: json_schema.MustNew(
			json_schema.RefDefinition("condition"),
		),
	}},
	Description: "A condition that negates the provided condition",
}

func (notCondition) compileCondition(ctx compilerContext) (condition, error) {
	p, decodeErr := notConditionSpec.Parse(ctx.jsonSpec)
	if decodeErr != nil {
		return condition{}, ctx.Error(decodeErr)
	}

	cond, cErr := compileCondition(ctx.child("not", p))
	if cErr != nil {
		return condition{}, ctx.Error(cErr)
	}

	return condition{
		check: func(ctx executionContext) (bool, error) {
			result, err := cond.check(ctx)
			return !result, err
		},
	}, nil
}

func (notCondition) JSONSchema() json_schema.JSONSchema {
	return notConditionSpec.JSONSchema()
}
