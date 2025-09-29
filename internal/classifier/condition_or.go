package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const orName = "or"

type orCondition struct{}

func (orCondition) name() string {
	return orName
}

var orConditionSpec = json_spec.SingleKeyValue[[]any]{
	Key: orName,
	ValueSpec: json_spec.MustSucceed[[]any]{json_spec.List[any]{
		ItemSpec: json_spec.Generic[any]{
			Schema: json_schema.MustNew(
				json_schema.RefDefinition("condition"),
			),
		},
		Description: "A condition that is satisfied if any of the conditions in a list are satisfied",
	}},
}

func (orCondition) compileCondition(ctx compilerContext) (condition, error) {
	rawConds, err := orConditionSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return condition{}, err
	}

	conds := make([]condition, len(rawConds))

	for i, rawCond := range rawConds {
		cond, err := compileCondition(ctx.child(json_spec.NumericPathPart(i), rawCond))
		if err != nil {
			return condition{}, err
		}

		conds[i] = cond
	}

	return condition{func(ctx executionContext) (bool, error) {
		for _, c := range conds {
			if result, err := c.check(ctx); err != nil {
				return false, err
			} else if result {
				return true, nil
			}
		}
		return false, nil
	}}, nil
}

func (orCondition) JSONSchema() json_schema.JSONSchema {
	return orConditionSpec.JSONSchema()
}
