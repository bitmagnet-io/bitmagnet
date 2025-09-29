package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

type ifElseAction struct{}

const ifElseName = "if_else"

func (ifElseAction) name() string {
	return ifElseName
}

type ifElse struct {
	Condition  any
	IfAction   any
	ElseAction any
}

var ifElseSpec = json_spec.SingleKeyValue[ifElse]{
	Key: ifElseName,
	ValueSpec: json_spec.MustSucceed[ifElse]{Typed: json_spec.Struct[ifElse]{
		Schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeObject),
			json_schema.Properties(map[string]json_schema.JSONSchema{
				"condition": json_schema.MustNew(
					json_schema.RefDefinition("condition"),
				),
				"if_action": json_schema.MustNew(
					json_schema.RefDefinition("action"),
				),
				"else_action": json_schema.MustNew(
					json_schema.RefDefinition("action"),
				),
			}),
			json_schema.Required(json_schema.RequiredFields{"condition"}),
			json_schema.AdditionalPropertiesFalse(),
		),
	}},
	Description: "Execute an action based on a condition",
}

func (ifElseAction) compile(ctx compilerContext) (action, error) {
	p, decodeErr := ifElseSpec.Parse(ctx.jsonSpec)
	if decodeErr != nil {
		return action{}, ctx.Error(decodeErr)
	}

	cond, cErr := compileCondition(ctx.child("condition", p.Condition))
	if cErr != nil {
		return action{}, ctx.Error(cErr)
	}

	var ifAction, elseAction action

	if p.IfAction != nil {
		pIfAction, ifErr := compileAction(ctx.child("if_action", p.IfAction))
		if ifErr != nil {
			return action{}, ctx.Error(ifErr)
		}

		ifAction = pIfAction
	}

	if p.ElseAction != nil {
		pElseAction, err := compileAction(ctx.child("else_action", p.ElseAction))
		if err != nil {
			return action{}, ctx.Error(err)
		}

		elseAction = pElseAction
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			if result, err := cond.check(ctx); err != nil {
				return classification.Result{}, err
			} else if result {
				if ifAction.run != nil {
					return ifAction.run(ctx)
				}
			} else {
				if elseAction.run != nil {
					return elseAction.run(ctx)
				}
			}
			return ctx.result, nil
		},
	}, nil
}

func (ifElseAction) JSONSchema() json_schema.JSONSchema {
	return ifElseSpec.JSONSchema()
}
