package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"

type ifElseAction struct{}

const ifElseName = "if_else"

func (ifElseAction) name() string {
	return ifElseName
}

type ifElsePayload struct {
	Condition  any
	IfAction   any
	ElseAction any
}

var ifElsePayloadSpec = payloadSingleKeyValue[ifElsePayload]{
	key: ifElseName,
	valueSpec: payloadMustSucceed[ifElsePayload]{payloadStruct[ifElsePayload]{
		jsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"condition": map[string]any{
					"$ref": "#/definitions/condition",
				},
				"if_action": map[string]any{
					"$ref": "#/definitions/action",
				},
				"else_action": map[string]any{
					"$ref": "#/definitions/action",
				},
			},
			"required":             []string{"condition"},
			"additionalProperties": false,
		},
	}},
	description: "Execute an action based on a condition",
}

func (ifElseAction) compileAction(ctx compilerContext) (action, error) {
	p, decodeErr := ifElsePayloadSpec.Unmarshal(ctx)
	if decodeErr != nil {
		return action{}, ctx.error(decodeErr)
	}

	cond, cErr := ctx.compileCondition(ctx.child("condition", p.Condition))
	if cErr != nil {
		return action{}, ctx.error(cErr)
	}

	var ifAction, elseAction action

	if p.IfAction != nil {
		pIfAction, ifErr := ctx.compileAction(ctx.child("if_action", p.IfAction))
		if ifErr != nil {
			return action{}, ctx.error(ifErr)
		}

		ifAction = pIfAction
	}

	if p.ElseAction != nil {
		pElseAction, err := ctx.compileAction(ctx.child("else_action", p.ElseAction))
		if err != nil {
			return action{}, ctx.error(err)
		}

		elseAction = pElseAction
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			ctx.logger.Info()
			logger:=ctx.logger
			ctx.logger=logger.Named("if_else.condition")
			result, err := cond.check(ctx)
			ctx.logger=logger
			if err != nil {
				ctx.logger.Info("error evaluating condition")
				return classification.Result{}, err
			} else if result {
				ctx.logger.Named("if_else.if_action").Info()
				if ifAction.run != nil {
					return ifAction.run(ctx)
				}
			} else {
				ctx.logger.Named("if_else.else_action").Info()
				if elseAction.run != nil {
					return elseAction.run(ctx)
				}
			}
			return ctx.result, nil
		},
	}, nil
}

func (ifElseAction) JSONSchema() JSONSchema {
	return ifElsePayloadSpec.JSONSchema()
}
