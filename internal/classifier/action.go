package classifier

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

func actions(defs ...actionDefinition) feature {
	return func(c *features) {
		c.actions = append(c.actions, defs...)
	}
}

type actionCompiler interface {
	compile(compilerContext) (action, error)
}

type actionDefinition interface {
	json_spec.HasJSONSchema
	actionCompiler
	name() string
}

func compileAction(ctx compilerContext) (action, error) {
	var (
		rawActions []any
		actions    []action
		errs       []error
	)

	isArray := false

	if s, ok := ctx.Source.([]any); ok {
		rawActions = s
		isArray = true
	} else {
		rawActions = []any{ctx.Source}
	}

outer:
	for i, rawAction := range rawActions {
		actionCtx := ctx
		if isArray {
			actionCtx = ctx.child(json_spec.NumericPathPart(i), rawAction)
		}
		for _, def := range ctx.actions {
			a, err := def.compile(actionCtx.child(def.name(), rawAction))
			if err == nil {
				actions = append(actions, a)
				continue outer
			}
			if json_spec.AsFatalCompilerError(err) != nil {
				return action{}, err
			}
		}
		errs = append(errs, fmt.Errorf("no action matched: %v", ctx.Source))
	}

	if len(errs) > 0 {
		return action{}, errors.Join(errs...)
	}

	return action{func(ctx executionContext) (classification.Result, error) {
		for _, a := range actions {
			result, err := a.run(ctx)
			if err != nil {
				return classification.Result{}, err
			}
			ctx = ctx.withResult(result)
		}
		return ctx.result, nil
	}}, errors.Join(errs...)
}

type action struct {
	run func(executionContext) (classification.Result, error)
}
