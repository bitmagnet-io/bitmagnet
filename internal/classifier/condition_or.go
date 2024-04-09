package classifier

const orName = "or"

type orCondition struct{}

func (orCondition) name() string {
	return orName
}

var orConditionSpec = payloadSingleKeyValue[[]any]{
	key: orName,
	valueSpec: payloadMustSucceed[[]any]{payloadList[any]{payloadGeneric[any]{
		jsonSchema: map[string]any{
			"$ref": "#/definitions/condition",
		},
	}}},
}

func (orCondition) compileCondition(ctx compilerContext) (condition, error) {
	rawConds, err := orConditionSpec.Unmarshal(ctx)
	if err != nil {
		return condition{}, err
	}
	conds := make([]condition, len(rawConds))
	for i, rawCond := range rawConds {
		cond, err := ctx.compileCondition(ctx.child(numericPathPart(i), rawCond))
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

func (orCondition) JsonSchema() JsonSchema {
	return orConditionSpec.JsonSchema()
}
