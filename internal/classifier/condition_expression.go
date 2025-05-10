package classifier

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
)

const expressionName = "expression"

type expressionCondition struct{}

var celProgramPayload = payloadTransformer[string, cel.Program]{
	spec: payloadGeneric[string]{
		jsonSchema: JSONSchema{
			"type":        "string",
			"minLength":   1,
			"description": "A CEL expression describing a condition",
		},
	},
	transform: func(s string, ctx compilerContext) (cel.Program, error) {
		ast, issues := ctx.celEnv.Compile(s)
		if issues != nil && issues.Err() != nil {
			return nil, ctx.error(fmt.Errorf("type-check error: %w", issues.Err()))
		}
		if !reflect.DeepEqual(ast.OutputType(), cel.BoolType) {
			return nil, ctx.error(
				fmt.Errorf("got %v, wanted %v output type", ast.OutputType(), cel.BoolType),
			)
		}
		prg, prgErr := ctx.celEnv.Program(ast,
			cel.EvalOptions(cel.OptOptimize),
		)
		if prgErr != nil {
			return nil, ctx.error(fmt.Errorf("program construction error: %w", prgErr))
		}
		return prg, nil
	},
}

var expressionConditionPayload = payloadUnion[cel.Program]{
	oneOf: []TypedPayload[cel.Program]{
		payloadSingleKeyValue[cel.Program]{
			key:       expressionName,
			valueSpec: payloadMustSucceed[cel.Program]{celProgramPayload},
		},
		payloadMustSucceed[cel.Program]{celProgramPayload},
	},
}

func (expressionCondition) name() string {
	return expressionName
}

func (expressionCondition) compileCondition(ctx compilerContext) (condition, error) {
	prg, err := expressionConditionPayload.Unmarshal(ctx)
	if err != nil {
		return condition{}, ctx.error(err)
	}

	return condition{
		check: func(ctx executionContext) (bool, error) {
			vars := map[string]any{
				"torrent": ctx.torrentPb,
				"result":  ctx.resultPb,
			}
			for k, v := range ctx.flags {
				vars["flags."+k] = v
			}
			result, _, err := prg.ContextEval(ctx.Context, vars)
			if err != nil {
				return false, err
			}
			bl, ok := result.Value().(bool)
			if !ok {
				return false, errors.New("not bool")
			}
			return bl, nil
		},
	}, nil
}

func (expressionCondition) JSONSchema() JSONSchema {
	return expressionConditionPayload.JSONSchema()
}
