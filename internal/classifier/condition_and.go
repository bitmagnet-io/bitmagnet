package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const andName = "and"

type andCondition struct{}

func (andCondition) name() string {
	return andName
}

var andConditionSpec = json_spec.SingleKeyValue[[]any]{
	Key: andName,
	ValueSpec: json_spec.MustSucceed[[]any]{json_spec.List[any]{
		ItemSpec: json_spec.Generic[any]{
			Schema: json_schema.MustNew(
				json_schema.RefDefinition("condition"),
			),
		},
		Description: "A condition that is satisfied if all conditions in a list are satisfied",
	}},
}

func (andCondition) compileCondition(ctx compilerContext) (condition, error) {
	payload, err := andConditionSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return condition{}, ctx.Error(err)
	}

	conds := make([]condition, len(payload))

	for i, rawCond := range payload {
		cond, err := compileCondition(ctx.child(json_spec.NumericPathPart(i), rawCond))
		if err != nil {
			return condition{}, ctx.Fatal(err)
		}

		conds[i] = cond
	}

	return condition{
		check: func(ctx executionContext) (bool, error) {
			for _, c := range conds {
				if result, err := c.check(ctx); err != nil {
					return false, err
				} else if !result {
					return false, nil
				}
			}
			return true, nil
		},
	}, nil
}

func (andCondition) JSONSchema() json_schema.JSONSchema {
	return andConditionSpec.JSONSchema()
}
