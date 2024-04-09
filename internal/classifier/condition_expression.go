package classifier

import (
	"errors"
	"fmt"
	"github.com/google/cel-go/cel"
	"reflect"
)

const expressionName = "expression"

type expressionCondition struct{}

var celProgramPayload = payloadTransformer[string, cel.Program]{
	spec: payloadGeneric[string]{
		jsonSchema: JsonSchema{
			"type":      "string",
			"minLength": 1,
		},
	},
	transform: func(s string, ctx compilerContext) (cel.Program, error) {
		ast, issues := ctx.celEnv.Compile(s)
		if issues != nil && issues.Err() != nil {
			return nil, ctx.error(fmt.Errorf("type-check error: %w", issues.Err()))
		}
		if !reflect.DeepEqual(ast.OutputType(), cel.BoolType) {
			return nil, ctx.error(fmt.Errorf("got %v, wanted %v output type", ast.OutputType(), cel.BoolType))
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

func (c expressionCondition) name() string {
	return expressionName
}

func (c expressionCondition) compileCondition(ctx compilerContext) (condition, error) {
	prg, err := expressionConditionPayload.Unmarshal(ctx)
	if err != nil {
		return condition{}, ctx.error(err)
	}
	return condition{
		check: func(ctx executionContext) (bool, error) {
			result, _, err := prg.Eval(map[string]any{
				"torrent": ctx.torrentPb,
				"result":  ctx.resultPb,
			})
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

func (c expressionCondition) JsonSchema() JsonSchema {
	return expressionConditionPayload.JsonSchema()
}
