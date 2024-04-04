package classifier

const andName = "and"

type andCondition struct{}

func (andCondition) name() string {
	return andName
}

var andConditionActionPayloadSpec = payloadSingleKeyValue[[]any]{
	key: andName,
	valueSpec: payloadMustSucceed[[]any]{payloadList[any]{payloadGeneric[any]{
		jsonSchema: map[string]any{
			"type": "any",
		},
	}}},
}

func (andCondition) compileCondition(ctx compilerContext) (condition, error) {
	payload, err := andConditionActionPayloadSpec.Unmarshal(ctx)
	if err != nil {
		return condition{}, ctx.error(err)
	}
	conds := make([]condition, len(payload))
	for i, rawCond := range payload {
		cond, err := ctx.compileCondition(ctx.child(numericPathPart(i), rawCond))
		if err != nil {
			return condition{}, ctx.fatal(err)
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
