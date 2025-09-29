package classifier

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/google/cel-go/cel"
)

const expressionName = "expression"

type expressionCondition struct{}

func celProgramPayload(env *cel.Env) json_spec.Transformer[string, cel.Program] {
	return json_spec.Transformer[string, cel.Program]{
		Typed: json_spec.Generic[string]{
			Schema: json_schema.MustNew(
				json_schema.Typed(json_schema.TypeString),
				json_schema.MinLength(1),
				json_schema.Description("A CEL expression describing a condition"),
			),
		},
		Transform: func(s string, ctx json_spec.ParseContext) (cel.Program, error) {
			ast, issues := env.Compile(s)
			if issues != nil && issues.Err() != nil {
				return nil, ctx.Error(fmt.Errorf("type-check error: %w", issues.Err()))
			}
			if !reflect.DeepEqual(ast.OutputType(), cel.BoolType) {
				return nil, ctx.Error(
					fmt.Errorf("got %v, wanted %v output type", ast.OutputType(), cel.BoolType),
				)
			}
			prg, prgErr := env.Program(ast,
				cel.EvalOptions(cel.OptOptimize),
			)
			if prgErr != nil {
				return nil, ctx.Error(fmt.Errorf("program construction error: %w", prgErr))
			}
			return prg, nil
		},
	}
}

func expressionConditionSpec(env *cel.Env) json_spec.Union[cel.Program] {
	valueSpec := json_spec.MustSucceed[cel.Program]{celProgramPayload(env)}

	return json_spec.Union[cel.Program]{
		OneOf: []json_spec.Typed[cel.Program]{
			json_spec.SingleKeyValue[cel.Program]{
				Key:       expressionName,
				ValueSpec: valueSpec,
			},
			valueSpec,
		},
	}
}

func (expressionCondition) name() string {
	return expressionName
}

func (expressionCondition) compileCondition(ctx compilerContext) (condition, error) {
	prg, err := expressionConditionSpec(ctx.celEnv).Parse(ctx.jsonSpec)
	if err != nil {
		return condition{}, ctx.Error(err)
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

func (expressionCondition) JSONSchema() json_schema.JSONSchema {
	return expressionConditionSpec(nil).JSONSchema()
}
