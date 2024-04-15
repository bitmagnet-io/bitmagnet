package classifier

const notName = "not"

type notCondition struct{}

func (notCondition) name() string {
	return notName
}

var notConditionPayloadSpec = payloadSingleKeyValue[any]{
	key: notName,
	valueSpec: payloadMustSucceed[any]{payloadGeneric[any]{
		jsonSchema: map[string]any{
			"$ref": "#/definitions/condition",
		},
	}},
}

func (notCondition) compileCondition(ctx compilerContext) (condition, error) {
	p, decodeErr := notConditionPayloadSpec.Unmarshal(ctx)
	if decodeErr != nil {
		return condition{}, ctx.error(decodeErr)
	}
	cond, cErr := ctx.compileCondition(ctx.child("not", p))
	if cErr != nil {
		return condition{}, ctx.error(cErr)
	}
	return condition{
		check: func(ctx executionContext) (bool, error) {
			result, err := cond.check(ctx)
			return !result, err
		},
	}, nil
}

func (notCondition) JsonSchema() JsonSchema {
	return notConditionPayloadSpec.JsonSchema()
}
